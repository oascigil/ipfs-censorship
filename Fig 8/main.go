package main

import (
    "context"
//     "flag"
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
    //rand.Seed(time.Now().UnixNano())
    rand.Seed(0)
}
