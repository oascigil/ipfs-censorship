package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"path/filepath"

	"github.com/ipfs/go-cid"
	// "github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	// kb "github.com/libp2p/go-libp2p-kbucket"
	"github.com/multiformats/go-multiaddr"
	rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
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
	TargetCID       string
	SybilIDs        []string
	ClosestPeers    []string
	PercentEclipsed float32
	ProviderPeerID  string
	IsEclipsed      bool
	UpdatedPeersUnmit	[]string
	ProviderDiscoveredPeersUnmit []string
	ContactedPeersUnmit	[]string
	RespondedPeersUnmit []string
	RegionSize int
	IsMitigated	 	bool
	UpdatedPeers	[]string
	ProviderDiscoveredPeers []string
	ContactedPeers	[]string
	RespondedPeers	[]string
	NumLookups	int
	Counts          []int
	Netsize         float64
	Threshold       float64
	KL              float64
	Detection       bool
}

type Error struct {
	Error string
}

func appendErrorJSON(writeJSON *[]byte, err string)  {
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

const(
	clientDHT = false
	serverDHT = true
)

func main() {
	config := Config{}
	Init()

	var numOfEclipses int // Number of CIDs to eclipse for each client
	var numOfLocations int // Number of clients from which to query the eclipsed CID
	var numOfSybils int // Number of Sybils to generate in each attack
	var regionSize int // Size of query region in mitigation, specified as expected number of honest peers in the region (20 is the value used in all experiments)
	var outpath string // path to store output of the experiment
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
	outpath = filepath.Join(outpath, "/sybil" + fmt.Sprint(numOfSybils) + "Combined")
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
			fmt.Println("Error marshaling JSON:", err)
		}
		writeJSON = append(writeJSON, locationJSON...)
		writeJSON = append(writeJSON, '\n')

		Init() // So that the same CIDs are eclipsed for all locs

		for eclipses := 0; eclipses < numOfEclipses; eclipses++ {
			fmt.Println("Loc", locs+1, "eclipse", eclipses+1)
			// update the provided file
			err = appendTextToFile(*fileName, "foo-bar-baz\n")
			if err != nil {
				log.Println(err)
			}
			// Get the CID of the content without providing the content
			cid_bytes := getCIDWithoutAdding(*fileName)
			targetCID := string(cid_bytes)
			targetCID = strings.ReplaceAll(targetCID, "\n", "")
			fmt.Printf("Target CID: '%s'\n", targetCID)

			out := getCurrentClosest(targetCID)
			temp := strings.Split(strings.TrimSpace(out), "\n")
			if out == "" {
				log.Println("Could not get closest nodes in DHT...")
				eclipses -= 1 // Try again with a different CID
				continue
			}
			currentClosest := temp[0]
			peerIdList := temp[1:]
			fmt.Println("Current Closest:", currentClosest)
			fmt.Println("Peer List:", peerIdList)

			pkeylist, sybilcidlist, err := attackCID(targetCID, currentClosest, peerIdList, numOfSybils)
			if err != nil {
				fmt.Println(err)
				eclipses -= 1
				continue
			}
			sybils := launchSybils(pkeylist, sybilcidlist)
			fmt.Println("Sleeping for a minute after launching Sybils...")
			time.Sleep(60 * time.Second)
			percentEclipsed := getPercentEclipsed(targetCID, sybilcidlist) // percentage of Sybil nodes in the 20 closest peers to the target CID

			cid, err := cid.Decode(targetCID)
			if err != nil {
				log.Printf("Failed to decode cid %s", targetCID)
				appendErrorJSON(&writeJSON, fmt.Sprintf("Failed to decode cid %s", targetCID))
				continue
			}
			multih := cid.Hash()

			peers, err := dhtClient.GetClosestPeers(ctx, string(multih))
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}
			if len(peers) == 0 {
				appendErrorJSON(&writeJSON, "0 closest peers found")
				continue
			}
			peersString := make([]string, 0, len(peers))
			for _, pid := range peers {
				peersString = append(peersString, pid.String())
			}
			log.Printf("Peer info: %q", peers)
			fmt.Println("Closest peers:", peers)
			detResult, err := dhtClient.EclipseDetectionVerbose(ctx, multih, peers)
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			// Now provide the content
			fmt.Println("Providing content...")
			dhtClient.DisableMitigation()
			dhtProvider.DisableMitigation()
			err, discoveredPeers, updatedPeers, _ := dhtProvider.ProvideWithReturn(ctx, cid, true)
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			// Check if eclipse actually succeeded
			provAddrInfos, contactedPeers, respondedPeers, err := dhtClient.FindProvidersReturnOnPathNodes(ctx, cid)
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			if len(provAddrInfos) > 0 {
				fmt.Println("\tCID was not eclipsed")
			} else {
				fmt.Println("\tCID was eclipsed")
			}

			updatedString := make([]string, 0, len(updatedPeers))
			for _, pid := range updatedPeers {
				updatedString = append(updatedString, pid.String())
			}
			discoveredString := make([]string, 0, len(discoveredPeers))
			for _, pid := range discoveredPeers {
				discoveredString = append(discoveredString, pid.String())
			}
			contactedString := make([]string, 0, len(contactedPeers))
			for _, pid := range contactedPeers {
				contactedString = append(contactedString, pid.String())
			}
			respondedString := make([]string, 0, len(respondedPeers))
			for _, pid := range respondedPeers {
				respondedString = append(respondedString, pid.String())
			}

			// Provide with mitigation

			dhtClient.EnableMitigation()
			dhtProvider.EnableMitigation()
			fmt.Println("Providing content with mitigation...")
			err, discoveredPeersMit, updatedPeersMit, numLookupsMit := dhtProvider.ProvideWithReturn(ctx, cid, true)
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			// Check if mitigation succeeded
			provAddrInfosMit, contactedPeersMit, respondedPeersMit, err := dhtClient.FindProvidersReturnOnPathNodes(ctx, cid)
			if err != nil {
				appendErrorJSON(&writeJSON, err.Error())
				continue
			}

			if len(provAddrInfosMit) > 0 {
				fmt.Println("\tAttack was mitigated")
			} else {
				fmt.Println("\tAttack was not mitigated")
			}

			updatedStringMit := make([]string, 0, len(updatedPeersMit))
			for _, pid := range updatedPeersMit {
				updatedStringMit = append(updatedStringMit, pid.String())
			}
			discoveredStringMit := make([]string, 0, len(discoveredPeersMit))
			for _, pid := range discoveredPeersMit {
				discoveredStringMit = append(discoveredStringMit, pid.String())
			}
			contactedStringMit := make([]string, 0, len(contactedPeersMit))
			for _, pid := range contactedPeersMit {
				contactedStringMit = append(contactedStringMit, pid.String())
			}
			respondedStringMit := make([]string, 0, len(respondedPeersMit))
			for _, pid := range respondedPeersMit {
				respondedStringMit = append(respondedStringMit, pid.String())
			}

			// intersection of updatedPeers and contactedPeers
			// remove Sybils
			// Assuming that there are no duplicates in these sets
			// nonSybilIntersectionMit := 0
			// for _, contacted := range contactedPeersMit {
			// 	for _, discovered := range discoveredPeersMit {
			// 		if contacted == discovered {
			// 			contactedString := contacted.String()
			// 			isSybil := false
			// 			for _, sybilCID := range allsybilcidlist {
			// 				if sybilCID == contactedString {
			// 					isSybil = true
			// 					break
			// 				}
			// 			}
			// 			if isSybil == false {
			// 				nonSybilIntersectionMit += 1
			// 			}
			// 		}
			// 	}
			// }
			// numSybilsUpdatedMit := 0
			// for _, discovered := range(discoveredPeersMit) {
			// 	discoveredString := discovered.String()
			// 	for _, sybilID := range(sybilcidlist) {
			// 		if (updatedString == sybilID) {
			// 			numSybilsUpdatedMit += 1
			// 		}
			// 	}
			// }

			experiment := Experiment{
				TargetCID:       targetCID,
				SybilIDs:        sybilcidlist,
				ClosestPeers:    peersString,
				PercentEclipsed: percentEclipsed,
				ProviderPeerID:  providerid,
				IsEclipsed:      len(provAddrInfos) == 0,
				UpdatedPeersUnmit: updatedString,
				ProviderDiscoveredPeersUnmit: discoveredString,
				ContactedPeersUnmit: contactedString,
				RespondedPeersUnmit: respondedString,
				RegionSize: regionSize,
				IsMitigated:	 len(provAddrInfosMit) > 0,
				ContactedPeers: 	contactedStringMit,
				UpdatedPeers:   	updatedStringMit,
				ProviderDiscoveredPeers: discoveredStringMit,
				RespondedPeers: respondedStringMit,
				NumLookups:   numLookupsMit,
				Counts:          detResult.Counts,
				Netsize:         detResult.Netsize,
				Threshold:       detResult.Threshold,
				KL:              detResult.KL,
				Detection:       detResult.Detection,
			}
			experimentJSON, err := json.Marshal(experiment)
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
			}
			writeJSON = append(writeJSON, experimentJSON...)
			writeJSON = append(writeJSON, '\n')

			killSybils(sybils)
			dhtProvider.Close()
			fmt.Println("Sleeping for ten seconds after killing Sybils...")
			time.Sleep(10 * time.Second)
		}
		dhtClient.Close()
		os.WriteFile(outpath+"/"+clientid+".json", writeJSON, 0644)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Finished test for clientID", clientid, ", written to file.")
		}
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

	// XXX do we need the code below (we can probably do: return dht here)

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
