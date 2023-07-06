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

	gocid "github.com/ipfs/go-cid"
	// "github.com/libp2p/go-libp2p-core/peer"
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

type Experiment struct {
	ClientPeerID          string
	ProviderPeerID        string
	NumSybils             int
	RegionSize            int
	ProvideLatencyMs      []int32
	FindProvsLatencyMs    []int32
	ProvideMitLatencyMs   []int32
	FindProvsMitLatencyMs []int32
	NumProvsFound         []int
	NumProvsFoundMit      []int
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

	var numOfSybils int // Number of Sybils to generate in each attack
	var numClients int  // Number of clients from which to query the eclipsed CID
	var numRuns int     // Number of CIDs to test for each client
	var regionSize int  // Size of query region in mitigation, specified as expected number of honest peers in the region (20 is the value used in all experiments)
	var outpath string  // path to store output of the experiment
	// var toMitigate bool // flag to measure latency with mitigation
	flag.IntVar(&numOfSybils, "sybils", 20, "Number of Sybils to generate")
	flag.IntVar(&numClients, "clients", 5, "Number of clients to generate")
	flag.IntVar(&numRuns, "runs", 1, "Number of runs to generate")
	flag.IntVar(&regionSize, "region", 20, "Region size for mitigation")
	flag.StringVar(&outpath, "outpath", "", "Path for output logs")
	fileName := flag.String("fileName", "./data/cat.txt", "File to Provide")
	// flag.BoolVar(&toMitigate, "mitigation", false, "Set this option to enable mitigation")
	flag.Parse()

	if outpath == "" {
		fmt.Println("Path for output logs not specified. Use -outfile <path>")
		os.Exit(1)
	}
	// TODO: adjust outpath based on from where python reads data
	err := os.MkdirAll(outpath, os.ModePerm)
	if err != nil {
		log.Println("Could not create directory", outpath)
		log.Println(err)
		os.Exit(1)
	}

	ctx, _ := context.WithCancel(context.Background())

	// Start the experiment:
	fmt.Printf("Experiments with %d Sybils...\n\n", numOfSybils)

	writeJSON := make([]byte, 0)

	for client := 0; client < numClients; client++ {
		fmt.Println("Client number ", client+1)
		// Create a client DHT node in a random location
		dhtClient, clientid := newDHTNode(config, ctx, clientDHT)
		dhtProvider, providerid := newDHTNode(config, ctx, serverDHT)
		dhtClient.SetProvideRegionSize(regionSize)
		dhtProvider.SetProvideRegionSize(regionSize)

		Init() // So that the same CIDs are eclipsed for all locs

		var provideLat []int32
		var provideMitLat []int32
		var findprovsLat []int32
		var findprovsMitLat []int32
		var numProvsFound []int
		var numProvsFoundMit []int

		for runs := 0; runs < numRuns; runs++ {
			fmt.Println("Client", client+1, ", run", runs+1)
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

			out := getCurrentClosest(targetCID)
			temp := strings.Split(strings.TrimSpace(out), "\n")
			if out == "" {
				fmt.Println("Could not get closest nodes in DHT, skipping...")
				runs -= 1
				continue
			}
			currentClosest := temp[0]
			peerIdList := temp[1:]
			fmt.Println("Current Closest:", currentClosest)
			fmt.Println("Peer List:", peerIdList)

			pkeylist, sybilcidlist, err := attackCID(targetCID, currentClosest, peerIdList, numOfSybils)
			if err != nil {
				fmt.Println("Took too long to generate Sybil IDs, skipping...")
				fmt.Println(err)
				runs -= 1
				continue
			}
			sybils := launchSybils(pkeylist, sybilcidlist)
			if numOfSybils > 0 {
				fmt.Println("Sleeping for a minute after launching Sybils...")
				time.Sleep(60 * time.Second)
			}

			cid, err := gocid.Decode(targetCID)
			if err != nil {
				fmt.Printf("Failed to decode cid %s\n", targetCID)
				runs -= 1
				continue
			}

			if numOfSybils == 0 {
				// Run default operations only when there is no attack
				// Now provide the content
				fmt.Println("Providing content...")
				dhtClient.DisableMitigation()
				dhtProvider.DisableMitigation()
				startTime := time.Now()
				err = dhtProvider.Provide(ctx, cid, true)
				thisLat := int32(time.Since(startTime).Milliseconds())
				fmt.Println("Finished provide")
				fmt.Println("Provide latency:", thisLat)
				if err != nil {
					fmt.Println(err)
					continue
				}
				provideLat = append(provideLat, thisLat)

				// Find providers without mitigation
				startTime = time.Now()
				provs, err := dhtClient.FindProviders(ctx, cid)
				thisLat = int32(time.Since(startTime).Milliseconds())
				fmt.Println("Finished find providers")
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println("Find providers latency:", thisLat)
				findprovsLat = append(findprovsLat, thisLat)
				numProvsFound = append(numProvsFound, len(provs))
			}

			// Provide with mitigation
			dhtClient.EnableMitigation()
			dhtProvider.EnableMitigation()
			fmt.Println("Providing content with mitigation...")
			startTime := time.Now()
			err = dhtProvider.Provide(ctx, cid, true)
			thisLat := int32(time.Since(startTime).Milliseconds())
			fmt.Println("Finished provide with mitigation")
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Provide latency with mitigation:", thisLat)
			provideMitLat = append(provideMitLat, thisLat)

			// Find providers with mitigation
			startTime = time.Now()
			provs, err := dhtClient.FindProviders(ctx, cid)
			thisLat = int32(time.Since(startTime).Milliseconds())
			fmt.Println("Finished find providers with mitigation")
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Find providers latency with mitigation:", thisLat)
			findprovsMitLat = append(findprovsMitLat, thisLat)
			numProvsFoundMit = append(numProvsFoundMit, len(provs))

			killSybils(sybils)
			fmt.Println("Sleeping for ten seconds after killing Sybils...")
			time.Sleep(10 * time.Second)
		}
		dhtProvider.Close()
		dhtClient.Close()

		experiment := Experiment{
			ClientPeerID:       clientid,
			ProviderPeerID:     providerid,
			NumSybils:          numOfSybils,
			RegionSize:         regionSize,
			ProvideLatencyMs:   provideLat,
			FindProvsLatencyMs: findprovsLat,
			ProvideMitLatencyMs:   provideMitLat,
			FindProvsMitLatencyMs: findprovsMitLat,
			NumProvsFound:      numProvsFound,
			NumProvsFoundMit:   numProvsFoundMit,
		}
		experimentJSON, err := json.Marshal(experiment)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		writeJSON = append(writeJSON, experimentJSON...)
		writeJSON = append(writeJSON, '\n')
	}
	os.WriteFile(outpath+"/latency.json", writeJSON, 0644)
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
