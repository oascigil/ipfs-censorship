package main

import (
    "context"
	"encoding/json"
    "flag"
    "fmt"
    "log"
    "os"
    "strings"
    "time"
    "math/rand"
    "github.com/multiformats/go-multiaddr"
    dht "github.com/libp2p/go-libp2p-kad-dht"
    "github.com/ipfs/go-cid"
    "github.com/libp2p/go-libp2p-core/peer"
    rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

    
func attackSucces(output string)bool{
    if strings.Contains(output,"Providers: []"){
        fmt.Println("No providers found, successful attack")
        return true
        
    }else{
        fmt.Println("Providers found...")
        outputA := strings.Split(output, "\n\n")
        fmt.Println(outputA[1])
        return false
    }

}

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
    ContactedPeers  []string
    UpdatedPeers    []string
    RespondedPeers []string
    ProviderDiscoveredPeers []string
    PercentEclipsed float32
    ProviderPeerID  string
    IsEclipsed      bool
    NumContacted    int
    NumUpdated      int
    NumSybilsUpdated int
    NumIntersection int 
    SpecialProvideNumber int
    IsMitigated     bool
    NumSuccessfulPings int
    NumLookups      int
    Counts          []int
    Netsize         float64
    Threshold       float64
    KL              float64
    Detection       bool
}


func Init() {
    //rand.Seed(time.Now().UnixNano())
    rand.Seed(0)
}
    
const(
	clientDHT = false
	serverDHT = true
)

func main(){
    config := Config{}
    Init()   

    var numOfEclipses int
    var numOfLocations int
    var numOfSybils int
	var specialProvideNumber int
	var outpath string
    
    flag.IntVar(&numOfEclipses, "cids", 5, "Number of Eclipses")
    flag.IntVar(&numOfLocations, "clients", 5, "Number of Locations")
    flag.IntVar(&numOfSybils, "sybils", 45, "Maximum number of Sybils to experiment (should ideally be a multiple of 10)")
	flag.IntVar(&specialProvideNumber, "special", 20, "Special provide number")
	flag.StringVar(&outpath, "outpath", "./logs/", "Path for output logs")
    fileName := flag.String("fileName", "./data/cat.txt", "File to Provide")
    flag.Parse()
    
	if outpath == "" {
		fmt.Println("Path for output logs not specified. Use -outfile <path>")
		os.Exit(1)
	}
	err := os.MkdirAll(outpath, os.ModePerm)
	if err != nil {
		log.Println("Could not create directory", outpath)
		log.Println(err)
		os.Exit(1)
	}

    ctx,_ := context.WithCancel(context.Background())

	fmt.Printf("Experiments with %d Sybils %d region size ...\n\n", numOfSybils, specialProvideNumber)
            
    for locs:=0; locs < numOfLocations;locs++ {
		fmt.Println("Loc", locs+1)
		dhtClient, clientid := newDHTNode(config, ctx, clientDHT)
		dhtProvider, providerid := newDHTNode(config, ctx, serverDHT)
        dhtClient.EnableMitigation()
        dhtProvider.EnableMitigation()
		dhtClient.SetProvideRegionSize(specialProvideNumber)
		dhtProvider.SetProvideRegionSize(specialProvideNumber)
		location := Location{
			ClientPeerID: clientid,
		}
        if (numOfSybils == 0) {
            dhtClient.DisableMitigation()
            dhtProvider.DisableMitigation()
        }
		
        writeJSON := make([]byte, 0)
		locationJSON, err := json.Marshal(location)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		writeJSON = append(writeJSON, locationJSON...)
		writeJSON = append(writeJSON, '\n')

		Init() // So that the same CIDs are eclipsed for all locs
        
        for eclipses:=0; eclipses < numOfEclipses; eclipses++ {
			fmt.Println("Loc", locs+1, "eclipse", eclipses+1)
            
            // update the provided file
            err = appendTextToFile(*fileName, "foo\n")
            if err != nil {
                log.Println(err)
            }
            // Get the CID of the content without providing the content
            cid_bytes := getCIDWithoutAdding(*fileName)
            targetCID := string(cid_bytes)
            targetCID = strings.ReplaceAll(targetCID, "\n", "")
            fmt.Printf("Target CID: '%s'\n", targetCID)

            out := getCurrentClosest(targetCID)
            fmt.Printf("out: %q\n", out)
            temp := strings.Split(out,"\n")
            if out == ""{
                log.Println("Could not get closest nodes in DHT...")
            }
            currentClosest := temp[0]
            peerIdList := temp[1:]
			currentFarthest := temp[len(temp)-1]
            fmt.Println("Current Closest:",currentClosest)
            fmt.Println("Peer List:",peerIdList)
			fmt.Println("Current Farthest:", currentFarthest)
			
            pkeylist, sybilcidlist, err := attackCID(targetCID, currentClosest, peerIdList, numOfSybils)
			if err != nil {
                fmt.Println("Took too long to generate Sybil IDs, giving up")
				continue
			}
    
			allsybilcidlist := sybilcidlist
            sybils := launchSybils(pkeylist, sybilcidlist)
            if (numOfSybils > 0) {
    			fmt.Println("Sleeping for a minute after launching Sybils...")
	    		time.Sleep(60 * time.Second)
            }
            percentEclipsed := getPercentEclipsed(targetCID, sybilcidlist)
        
            cid, err := cid.Decode(targetCID)
            if (err != nil) {
                log.Printf("Failed to decode cid %s", targetCID)
                log.Fatal(err)
            }
			multih := cid.Hash()

			peers, err := dhtClient.GetClosestPeers(ctx, string(multih))
			if err != nil {
				log.Fatal(err)
			}
			if len(peers) == 0 {
				continue
			}
			peersString := make([]string, 0, len(peers))
			for _, pid := range peers {
				// c := []byte(kb.ConvertKey(string(pid)))
				peersString = append(peersString, fmt.Sprintf("%s",pid))
			}
			log.Printf("Peer info: %q", peers)
			fmt.Println("Closest peers:", peers)
			detResult, err := dhtClient.EclipseDetectionVerbose(ctx, multih, peers)
			if err != nil {
				log.Fatal(err)
			}

			// Now provide the content
            fmt.Println("Providing content with eclipse detection...")
            err, discoveredPeers, updatedPeers, numLookups := dhtProvider.ProvideWithReturn(ctx, cid, true)
            if (err != nil) {
                log.Fatal(err)
            }

            // Check if the client can mitigate eclipsing
            provAddrInfos, contactedPeers, respondedPeers, err := dhtClient.FindProvidersReturnOnPathNodes(ctx, cid)
            if (err != nil) {
                log.Fatal(err)
            }

            discoveredPeersList := make([]string, 0, len(discoveredPeers))
            for _, discovered := range(discoveredPeers) {
                discoveredString := discovered.String()
                discoveredPeersList = append(discoveredPeersList, discoveredString)
            }
            respondedList := make([]string, 0, len(respondedPeers))
            for _, responded := range(respondedPeers) {
                respondedString := responded.String()
                respondedList = append(respondedList, respondedString)
            }
            contactedPeersList := make([]string, 0, len(contactedPeers))
            // intersection of updatedPeers and contactedPeers
            // remove Sybils
            nonSybilIntersection := 0
            var intersection []peer.ID
            for _, contacted := range(contactedPeers) {
                contactedString := contacted.String()
                contactedPeersList = append(contactedPeersList, contactedString)
                for _, updated := range(updatedPeers) {
                    if contacted == updated {
                        //contactedString := contacted.String()
                        isSybil := false
                        for _, sybilCID := range(sybilcidlist) {
                            if sybilCID == contactedString {
                                isSybil = true
                                break
                            }
                        }
                        if isSybil == false {
                            nonSybilIntersection += 1
                            intersection = append(intersection, contacted)
                        }
                    }
                }
            }
            updatedPeersList := make([]string, 0, len(updatedPeers))
            numSybilsUpdated := 0
            for _, updatedPeer := range(updatedPeers) {
                updatedPeerString := updatedPeer.String()
                updatedPeersList = append(updatedPeersList, updatedPeerString)
                for _, sybilID := range(sybilcidlist) {
                    if (updatedPeerString == sybilID) {
                        numSybilsUpdated += 1
                    }
                }
            }

            if len(provAddrInfos) > 0 {
                fmt.Println("\tSuccess: attack was mitigated")
            } else {
                fmt.Println("\tFailure: unable to mitigate attack")
            }

            numPongs := 0
            if nonSybilIntersection > 0 { 
                // mitigation failed, ping the nodes in the non-Sybil intersection set
                for _, p := range(intersection) {
                    err := dhtClient.Ping(ctx, p)
                    if (err == nil) {
                        numPongs += 1
                    }
                }
            }

			experiment := Experiment{
				TargetCID:       targetCID,
				SybilIDs:        allsybilcidlist,
				ClosestPeers:    peersString,
                ContactedPeers:  contactedPeersList,
                UpdatedPeers:    updatedPeersList,
                RespondedPeers:  respondedList,
                ProviderDiscoveredPeers:    discoveredPeersList,
				PercentEclipsed: percentEclipsed,
				ProviderPeerID:  providerid,
				IsEclipsed:      percentEclipsed > 99, // as a proxy for IsEclipsed because FindProviders is not checked without mitigation in this experiment
				NumContacted:    len(contactedPeers),
				NumUpdated:      len(updatedPeers),
				NumSybilsUpdated: numSybilsUpdated,
				NumIntersection: nonSybilIntersection,
				SpecialProvideNumber: specialProvideNumber,
				IsMitigated:	 len(provAddrInfos) > 0,
				NumSuccessfulPings:   numPongs,
				NumLookups:   numLookups,
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
        }
		os.WriteFile(outpath+"/"+clientid+".json", writeJSON, 0644)
        dhtClient.Close()
        dhtProvider.Close()
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
    if (isServer == true) {
        dht, routedHost = NewDHTServer(ctx, h, config.DiscoveryPeers)
    }else {
        dht, routedHost = NewDHTClient(ctx, h, config.DiscoveryPeers)
    }

    // Build host multiaddress
    hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))

    // Now we can build a full multiaddress to reach this host
    // by encapsulating both addresses:
    // addr := routedHost.Addrs()[0]
    addrs := routedHost.Addrs()
    
    for _, addr := range addrs {
        log.Println(addr.Encapsulate(hostAddr))
    }

    return dht, h.ID().Pretty()
}
