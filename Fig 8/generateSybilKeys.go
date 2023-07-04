package main

import (
        "crypto/rand"
        "fmt"
        "github.com/libp2p/go-libp2p-core/crypto"
        "github.com/libp2p/go-libp2p-core/peer"
        kspace "github.com/libp2p/go-libp2p-kbucket/keyspace"
        mh "github.com/multiformats/go-multihash"
        "log"
)

func generateNewKey(currentClosest string, targetCIDKey kspace.Key) (string,kspace.Key,string,error) {

        currentClosestByte,_ := mh.FromB58String(currentClosest)
        currentClosestKey := kspace.XORKeySpace.Key(currentClosestByte)
        currentDistance := currentClosestKey.Distance(targetCIDKey)

        tries := 0
        for true {
                tries = tries+1
                privateKey, publicKey, err := crypto.GenerateEd25519Key(rand.Reader)
                if err != nil {
                        log.Println(err)
                        return "",kspace.Key{nil,nil,nil},"", err
                }

                peerID, err := peer.IDFromPublicKey(publicKey)
                if err != nil {
                        log.Println(err)
                        return "",kspace.Key{nil,nil,nil},"",err
                }

                prettyID := peerID.Pretty()

                if err != nil {
                        log.Println(err)
                        return "",kspace.Key{nil,nil,nil},"", err
                }
                newPeerID,_ := mh.FromB58String((prettyID))
                newPeerKey := kspace.XORKeySpace.Key(newPeerID)
                newDistance := newPeerKey.Distance(targetCIDKey)

                if err != nil {
                        fmt.Printf("Error when marshalling public key")
                }

                marshal_priv, err := crypto.MarshalPrivateKey(privateKey)
                if err != nil {
                        fmt.Printf("Error when marshalling private key")
                }

                if newDistance.Cmp(currentDistance) == -1 {

                        fmt.Printf("\tdistance: %d\n", newDistance)
                        fmt.Printf("\t")
                        fmt.Println("Sybil PeerID:", prettyID)
                        fmt.Printf("\t")
                        fmt.Println("Sybil Private Key:", crypto.ConfigEncodeKey(marshal_priv))
                        fmt.Println("Found in: ", tries, " tries")
                        return prettyID,newPeerKey,crypto.ConfigEncodeKey(marshal_priv), nil
                }

        }
        return "",kspace.Key{nil,nil,nil},"",nil
}
