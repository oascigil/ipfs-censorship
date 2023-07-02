package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
    "syscall"
    "math/rand"

    "github.com/libp2p/go-libp2p-core/host"
    "github.com/multiformats/go-multiaddr"
)

const base58HashLength = 46

type Config struct {
    Port           int
    Seed           int64
    DiscoveryPeers addrList
}

func Init() {
    //rand.Seed(time.Now().UnixNano())
    rand.Seed(0)
}

var base58chars = []rune("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
func RandBase58String(length int) *string {
    b := make([]rune, length)
    for i := range b {
        b[i] = base58chars[rand.Intn(len(base58chars))]
    }
    b[0] = 'Q'
    b[1] = 'm'
    b[2] = 'Q' // apparently, the 2nd character determines the length of cid and must be a certain character ('Y' also works)
    s := string(b)
    return &s
}

func main() {
    config := Config{}
    Init()

    flag.Int64Var(&config.Seed, "seed", 1, "Seed value for generating a PeerID, 0 is random")
    flag.Var(&config.DiscoveryPeers, "peer", "Peer multiaddress for peer discovery")
    flag.IntVar(&config.Port, "port", 0, "")
    privKeyStr := flag.String("privKey", "CAESQFAOFXLJF2EQl9kkd63+COLXMSpskmGN8IdbOJa1mAUND2tuv2w17HgfUoXWhiYjgUlgq3dfZKPOXyTsnbf67Tc=", "Private Key")
    contentID := flag.String("cid", "QmTz3oc4gdpRMKP2sdGUPZTAGRngqjsi99BPoztyP53JMM", "Content ID")
    flag.Parse()

    ctx, cancel := context.WithCancel(context.Background())

    h, err := NewHost(ctx, config.Seed, config.Port, *privKeyStr)
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

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	// addr := routedHost.Addrs()[0]
	addrs := routedHost.Addrs()
	log.Println("I can be reached at:")
	for _, addr := range addrs {
		log.Println(addr.Encapsulate(hostAddr))
	}

    
    peers, err := dht.GetClosestPeers(ctx, string(*contentID))
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Peer info: %q", peers)

    run(h, cancel)
}

func run(h host.Host, cancel func()) {
    c := make(chan os.Signal, 1)

    signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
    <-c

    fmt.Printf("\rExiting...\n")

    cancel()

    if err := h.Close(); err != nil {
        panic(err)
    }
    os.Exit(0)
}

type addrList []multiaddr.Multiaddr

func (al *addrList) String() string {
    strs := make([]string, len(*al))
    for i, addr := range *al {
        strs[i] = addr.String()
    }
    return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
    addr, err := multiaddr.NewMultiaddr(value)
    if err != nil {
        return err
    }
    *al = append(*al, addr)
    return nil
}
