package main

import (
	"context"
	"log"
	//"sync"

	"github.com/libp2p/go-libp2p-core/host"
	//"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
    ds "github.com/ipfs/go-datastore"
    dsync "github.com/ipfs/go-datastore/sync"
    rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
    //bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
)

func NewDHT(ctx context.Context, basicHost host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, *rhost.RoutedHost) {

    // Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
    dstore := dsync.MutexWrap(ds.NewMapDatastore())

    // Make the DHT
	dht := dht.NewDHT(ctx, basicHost, dstore)

    // Make the routed host
	routedHost := rhost.Wrap(basicHost, dht)

    // connect to the chosen ipfs nodes
	//err = bootstrapConnect(ctx, routedHost, bootstrapPeers)
	err := bootstrapConnect(ctx, routedHost, IPFS_PEERS)
	if err != nil {
        log.Fatal(err)
	}

    // Bootstrap the host
	err = dht.Bootstrap(ctx)
	if err != nil {
        log.Fatal(err)
	}

    return dht, routedHost

    /*
	var options []dht.Option

    bootstrapPeers = dht.DefaultBootstrapPeers

	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

    // Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
    dstore := dsync.MutexWrap(ds.NewMapDatastore())
	//kdht, err := dht.New(ctx, host, options...)
	//if err != nil {
	//	return nil, err
	//}
    //kdht := dht.NewDHTClient(ctx, host, options...)

    //Onur: using DHTClient instead 
    kdht := dht.NewDHTClient(ctx, host, dstore)

	err := kdht.Bootstrap(ctx)
    if err != nil {
		return nil, err
	}


    // Connect to bootstrap nodes
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return kdht, nil 
    */
}
