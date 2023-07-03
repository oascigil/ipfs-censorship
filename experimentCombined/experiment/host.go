package main

import (
	"context"
	"crypto/rand"
	"fmt"
    "log"
	"io"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	//"github.com/multiformats/go-multiaddr"
)

func NewHost(ctx context.Context, seed int64, listenPort int) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mrand.New(mrand.NewSource(seed))
	}

	privKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
        log.Fatal(err)
		return nil, err
	}

    opts := []libp2p.Option{
	libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort)),
	libp2p.Identity(privKey),
	libp2p.DefaultTransports,
	libp2p.DefaultMuxers,
	libp2p.DefaultSecurity,
	libp2p.NATPortMap(),
	}

	return libp2p.New(opts...)

	/*addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))

	return libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(priv),
	)*/
}
