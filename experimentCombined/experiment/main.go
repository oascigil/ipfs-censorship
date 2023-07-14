package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	gocid "github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	// kb "github.com/libp2p/go-libp2p-kbucket"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/multiformats/go-multiaddr"
)

type addrList []multiaddr.Multiaddr

type Config struct {
	Port           int
	Seed           int64
	DiscoveryPeers addrList
}

type Location struct {
	ClientPeerID string
}

type Experiment struct {
	TargetCID                    string
	SybilIDs                     []string
	ClosestPeers                 []string
	PercentEclipsed              float64
	ProviderPeerID               string
	IsEclipsed                   bool
	RegionSize                   int
	IsMitigated                  bool
	UpdatedPeers                 []string
	ProviderDiscoveredPeers      []string
	ContactedPeers               []string
	RespondedPeers               []string
	NumContacted                 int
	NumUpdated                   int
	NumSybilsUpdated             int
	NumIntersection              int
	NumSuccessfulPings           int
	NumLookups                   int
	Counts                       []int
	Netsize                      float64
	Threshold                    float64
	KL                           float64
	Detection                    bool
}

type Error struct {
	Error string
}

func appendErrorJSON(writeJSON *[]byte, err string) {
	errorJSON, err2 := json.Marshal(Error{Error: err})
	if err2 != nil {
		fmt.Println("Error marshaling JSON:", err2)
	}
	*writeJSON = append(*writeJSON, errorJSON...)
	*writeJSON = append(*writeJSON, '\n')
}

func Init() {
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(0) // So that the same CIDs are eclipsed for all locs
}

const (
	clientDHT = false
	serverDHT = true
)

func main() {
	config := Config{}
	Init()

	var numOfEclipses int  // Number of CIDs to eclipse for each client
	var numOfLocations int // Number of clients from which to query the eclipsed CID
	var numOfSybils int    // Number of Sybils to generate in each attack
	var regionSize int     // Size of query region in mitigation, specified as expected number of honest peers in the region (20 is the value used in all experiments)
	var outpath string     // path to store output of the experiment
	flag.IntVar(&numOfEclipses, "cids", 1, "Number of CIDs to eclipse")
	flag.IntVar(&numOfLocations, "clients", 1, "Number of clients to test from")
	flag.IntVar(&numOfSybils, "sybils", 45, "Number of Sybils to generate")
	flag.IntVar(&regionSize, "region", 20, "Region size for mitigation")
	flag.StringVar(&outpath, "outpath", "", "Path for output logs")
	fileName := flag.String("filename", "./data/cat.txt", "File to Provide")
	flag.Parse()

	if outpath == "" {
		fmt.Println("Path for output logs not specified. Use -outfile <path>")
		os.Exit(1)
	}
	outpath = filepath.Join(outpath, "/sybil"+fmt.Sprint(numOfSybils)+"Combined")
	err := os.MkdirAll(outpath, os.ModePerm)
	if err != nil {
		log.Println("Could not create directory", outpath)
		log.Println(err)
		os.Exit(1)
	}

	ctx, _ := context.WithCancel(context.Background())

	// Start the experiment:
	fmt.Printf("Experiments with %d Sybils...\n\n", numOfSybils)

	for locs := 0; locs < numOfLocations; locs++ {
		fmt.Println("Loc", locs+1)
		// Create a client DHT node in a random location
		dhtClient, clientid := newDHTNode(config, ctx, clientDHT)
		dhtProvider, providerid := newDHTNode(config, ctx, serverDHT)
		dhtClient.SetProvideRegionSize(regionSize)
		dhtProvider.SetProvideRegionSize(regionSize)
		location := Location{
			ClientPeerID: clientid,
		}

		writeJSON := make([]byte, 0)
		locationJSON, err := json.Marshal(location)
		if err != nil {
			log.Fatal("Error marshaling JSON:", err)
		}
		writeJSON = append(writeJSON, locationJSON...)
		writeJSON = append(writeJSON, '\n')

		Init() // So that the same CIDs are eclipsed for all locs

		for eclipses := 0; eclipses < numOfEclipses; eclipses++ {
			fmt.Println("Loc", locs+1, "eclipse", eclipses+1)
			// update the provided file
			err = appendTextToFile(*fileName, "foo-bar-baz\n")
			if err != nil {
				log.Fatal(err)
			}
			// Get the CID of the content without providing the content
			cid_bytes := getCIDWithoutAdding(*fileName)
			targetCID := string(cid_bytes)
			targetCID = strings.ReplaceAll(targetCID, "\n", "")
			fmt.Printf("Target CID: '%s'\n", targetCID)

			var sybilcidlist []string
			var sybils []*exec.Cmd
			if numOfSybils > 0 {
				out, err := getCurrentClosest(targetCID, 120 * time.Second)
				temp := strings.Split(strings.TrimSpace(out), "\n")
				if out == "" || err != nil {
					fmt.Println("Could not get closest nodes in DHT, skipping...")
					eclipses -= 1 // Try again with a different CID
					continue
				}
				currentClosest := temp[0]
				peerIdList := temp[1:]
				fmt.Println("Current Closest:", currentClosest)
				fmt.Println("Peer List:", peerIdList)

				var pkeylist []string
				pkeylist, sybilcidlist, err = attackCID(targetCID, currentClosest, peerIdList, numOfSybils)
				if err != nil {
					fmt.Println("Took too long to generate Sybil IDs, skipping...")
					fmt.Println(err)
					eclipses -= 1
					continue
				}
				sybils = launchSybils(pkeylist, sybilcidlist)
				fmt.Println("Sleeping for a minute after launching Sybils...")
				time.Sleep(60 * time.Second)
			}

			cid, err := gocid.Decode(targetCID)
			if err != nil {
				fmt.Printf("Failed to decode cid %s\n", targetCID)
				appendErrorJSON(&writeJSON, fmt.Sprintf("Failed to decode cid %s", targetCID))
				continue
			}
			multih := cid.Hash()

			ctxWithTimeout, _ := context.WithTimeout(ctx, 600*time.Second)
			peers, err := dhtClient.GetClosestPeers(ctxWithTimeout, string(multih))
			if err != nil {
				fmt.Println("Error getting closest peers:", err)
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}
			if len(peers) == 0 {
				fmt.Println("0 closest peers found")
				appendErrorJSON(&writeJSON, "0 closest peers found")
				continue
			}
			peersString := make([]string, 0, len(peers))
			for _, pid := range peers {
				peersString = append(peersString, pid.String())
			}
			fmt.Println("Closest peers:", peers)

			percentEclipsed := 0.0 // percentage of Sybil nodes in the 20 closest peers to the target CID
			for i := 0; i < len(peersString); i++ {
				for j := 0; j < len(sybilcidlist); j++ {
					if peersString[i] == sybilcidlist[j] {
						percentEclipsed += 1
					}
				}
			}
			percentEclipsed = percentEclipsed / 20 * 100

			ctxWithTimeout, _ = context.WithTimeout(ctx, 600*time.Second)
			detResult, err := dhtClient.EclipseDetectionVerbose(ctxWithTimeout, multih, peers)
			if err != nil {
				fmt.Println("Error in eclipse detection:", err)
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			// Now provide the content
			fmt.Println("Providing content...")
			dhtClient.DisableMitigation()
			dhtProvider.DisableMitigation()
			ctxWithTimeout, _ = context.WithTimeout(ctx, 600*time.Second)
			err, discoveredPeersUnmit, updatedPeersUnmit, _ := dhtProvider.ProvideWithReturn(ctxWithTimeout, cid, true)
			if err != nil {
				fmt.Println("Error providing content:", err)
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			// Check if attack succeeded
			ctxWithTimeout, _ = context.WithTimeout(ctx, 600*time.Second)
			provAddrInfosUnmit, contactedPeersUnmit, respondedPeersUnmit, err := dhtClient.FindProvidersReturnOnPathNodes(ctxWithTimeout, cid)
			if err != nil {
				fmt.Println("Error finding providers:", err)
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			if len(provAddrInfosUnmit) > 0 {
				fmt.Println("\tCID was not censored")
			} else {
				fmt.Println("\tCID was censored")
			}

			updatedStringUnmit := make([]string, 0, len(updatedPeersUnmit))
			for _, pid := range updatedPeersUnmit {
				updatedStringUnmit = append(updatedStringUnmit, pid.String())
			}
			discoveredStringUnmit := make([]string, 0, len(discoveredPeersUnmit))
			for _, pid := range discoveredPeersUnmit {
				discoveredStringUnmit = append(discoveredStringUnmit, pid.String())
			}
			contactedStringUnmit := make([]string, 0, len(contactedPeersUnmit))
			for _, pid := range contactedPeersUnmit {
				contactedStringUnmit = append(contactedStringUnmit, pid.String())
			}
			respondedStringUnmit := make([]string, 0, len(respondedPeersUnmit))
			for _, pid := range respondedPeersUnmit {
				respondedStringUnmit = append(respondedStringUnmit, pid.String())
			}

			var discoveredString, updatedString, contactedString, respondedString []string
			var updatedPeers, contactedPeers []peer.ID
			var numLookups int
			var provAddrInfos []peer.AddrInfo
			if numOfSybils > 0 {
				// Provide with mitigation
				dhtClient.EnableMitigation()
				dhtProvider.EnableMitigation()
				fmt.Println("Providing content with mitigation...")
				var discoveredPeers []peer.ID
				ctxWithTimeout, _ = context.WithTimeout(ctx, 600*time.Second)
				err, discoveredPeers, updatedPeers, numLookups = dhtProvider.ProvideWithReturn(ctxWithTimeout, cid, true)
				if err != nil {
					fmt.Println("Error providing content:", err)
					appendErrorJSON(&writeJSON, err.Error())
					continue
				}

				// Check if mitigation succeeded
				var respondedPeers []peer.ID
				ctxWithTimeout, _ = context.WithTimeout(ctx, 600*time.Second)
				provAddrInfos, contactedPeers, respondedPeers, err = dhtClient.FindProvidersReturnOnPathNodes(ctxWithTimeout, cid)
				if err != nil {
					fmt.Println("Error finding providers:", err)
					appendErrorJSON(&writeJSON, err.Error())
					continue
				}

				if len(provAddrInfos) > 0 {
					fmt.Println("\tSuccess: attack was mitigated")
				} else {
					fmt.Println("\tFailure: unable to mitigate attack")
				}

				discoveredString = make([]string, 0, len(discoveredPeers))
				for _, pid := range discoveredPeers {
					discoveredString = append(discoveredString, pid.String())
				}
				contactedString = make([]string, 0, len(contactedPeers))
				for _, pid := range contactedPeers {
					contactedString = append(contactedString, pid.String())
				}
				respondedString = make([]string, 0, len(respondedPeers))
				for _, pid := range respondedPeers {
					respondedString = append(respondedString, pid.String())
				}
				updatedString = make([]string, 0, len(updatedPeers))
				for _, updated := range updatedPeers {
					updatedString = append(updatedString, updated.String())
				}
			} else {
				// No sybils, so no mitigation
				discoveredString = discoveredStringUnmit
				updatedString = updatedStringUnmit
				contactedString = contactedStringUnmit
				respondedString = respondedStringUnmit
				contactedPeers = contactedPeersUnmit
				updatedPeers = updatedPeersUnmit
				provAddrInfos = provAddrInfosUnmit
				numLookups = 1
			}

			// intersection of updatedPeers and contactedPeers
			// remove Sybils
			// Assuming that there are no duplicates in these sets
			nonSybilIntersection := 0
			var intersection []peer.ID
			for _, contacted := range contactedPeers {
				for _, updated := range updatedPeers {
					if contacted == updated {
						isSybil := false
						for _, sybilCID := range sybilcidlist {
							if sybilCID == contacted.String() {
								isSybil = true
								break
							}
						}
						if !isSybil {
							nonSybilIntersection += 1
							intersection = append(intersection, contacted)
						}
					}
				}
			}
			numSybilsUpdated := 0
			for _, updated := range updatedString {
				for _, sybilID := range sybilcidlist {
					if updated == sybilID {
						numSybilsUpdated += 1
					}
				}
			}
			numPongs := 0
			if nonSybilIntersection > 0 {
				// if mitigation failed even with non-Sybil peers in the intersection, we want to ping the nodes in the non-Sybil intersection set to investigate
				for _, p := range intersection {
					err := dhtClient.Ping(ctx, p)
					if err == nil {
						numPongs += 1
					}
				}
			}

			experiment := Experiment{
				TargetCID:       targetCID,
				SybilIDs:        sybilcidlist,
				ClosestPeers:    peersString,
				PercentEclipsed: percentEclipsed,
				ProviderPeerID:  providerid,
				IsEclipsed:      len(provAddrInfosUnmit) == 0,
				RegionSize:              regionSize,
				IsMitigated:             len(provAddrInfos) > 0,
				ContactedPeers:          contactedString,
				UpdatedPeers:            updatedString,
				ProviderDiscoveredPeers: discoveredString,
				RespondedPeers:          respondedString,
				NumContacted:    len(contactedString),
				NumUpdated:      len(updatedString),
				NumSybilsUpdated: numSybilsUpdated,
				NumIntersection: nonSybilIntersection,
				NumLookups:              numLookups,
				NumSuccessfulPings:   numPongs,
				Counts:                  detResult.Counts,
				Netsize:                 detResult.Netsize,
				Threshold:               detResult.Threshold,
				KL:                      detResult.KL,
				Detection:               detResult.Detection,
			}
			experimentJSON, err := json.Marshal(experiment)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}
			writeJSON = append(writeJSON, experimentJSON...)
			writeJSON = append(writeJSON, '\n')
			if numOfSybils > 0 {
				killSybils(sybils)
				fmt.Println("Sleeping for ten seconds after killing Sybils...")
				time.Sleep(10 * time.Second)
			}
		}
		dhtClient.Close()
		dhtProvider.Close()
		os.WriteFile(outpath+"/"+clientid+".json", writeJSON, 0644)

		fmt.Println("Finished test for clientID", clientid, ", written to file.")
	}
}

func newDHTNode(config Config, ctx context.Context, isServer bool) (*dht.IpfsDHT, string) {

	h, err := NewHost(ctx, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Host ID: %s", h.ID().Pretty())
	log.Printf("Connect to me on:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
	}

	var dht *dht.IpfsDHT
	var routedHost *rhost.RoutedHost
	dht, routedHost = NewDHT(ctx, h, config.DiscoveryPeers, isServer)

	// Build host multiaddress
	hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addrs := routedHost.Addrs()

	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

	return dht, h.ID().Pretty()
}
