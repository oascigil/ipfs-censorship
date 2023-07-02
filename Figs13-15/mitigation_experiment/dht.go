package main

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
    ds "github.com/ipfs/go-datastore"
    dsync "github.com/ipfs/go-datastore/sync"
    rhost "github.com/libp2p/go-libp2p/p2p/host/routed"
)

func NewDHTServer(ctx context.Context, basicHost host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, *rhost.RoutedHost) {

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
}

func NewDHTClient(ctx context.Context, basicHost host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, *rhost.RoutedHost) {

    // Construct a datastore (needed by the DHT). This is just a simple, in-memory thread-safe datastore.
    dstore := dsync.MutexWrap(ds.NewMapDatastore())

    // Make the DHT
	dht := dht.NewDHTClient(ctx, basicHost, dstore)

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
}
