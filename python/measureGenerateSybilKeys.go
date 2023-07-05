package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	kspace "github.com/libp2p/go-libp2p-kbucket/keyspace"
	mh "github.com/multiformats/go-multihash"
)

func generateNewKey(targetCIDKey kspace.Key, networkSize int, useRSA bool) int {
	var maxDistance big.Int
	maxDistance.Exp(big.NewInt(2), big.NewInt(256), nil)

	threshold := new(big.Int).Div(&maxDistance, big.NewInt(int64(networkSize)))

	//fmt.Printf("\tmaxDistance: %d\n", maxDistance)
	//fmt.Printf("\tthreshold: %d\n", threshold)
	tries := 0
	for true {
		tries = tries + 1
		publicKey := crypto.PubKey(nil)
		err := error(nil)
		if useRSA {
			_, publicKey, err = crypto.GenerateRSAKeyPair(2048, rand.Reader)
		} else {
			_, publicKey, err = crypto.GenerateEd25519Key(rand.Reader)
		}

		if err != nil {
			log.Println(err)
			return -1
		}

		if err != nil {
			log.Println(err)
			return -1
		}

		peerID, err := peer.IDFromPublicKey(publicKey)
		if err != nil {
			log.Println(err)
			return -1
		}

		prettyID := peerID.Pretty()

		if err != nil {
			log.Println(err)
			return -1
		}
		newPeerID, _ := mh.FromB58String((prettyID))
		newPeerKey := kspace.XORKeySpace.Key(newPeerID)

		newDistance := newPeerKey.Distance(targetCIDKey)

		if err != nil {
			fmt.Printf("Error when marshalling public key")
		}

		/*marshal_priv, err := crypto.MarshalPrivateKey(privateKey)
		if err != nil {
			fmt.Printf("Error when marshalling private key")
		}*/

		if newDistance.Cmp(threshold) == -1 {

			/*fmt.Printf("\tdistance: %d\n", newDistance)
			fmt.Printf("\t")
			fmt.Println("Sybil PeerID:", prettyID)
			fmt.Printf("\t")
			fmt.Println("Sybil Private Key:", crypto.ConfigEncodeKey(marshal_priv))
			fmt.Println("Found in: ", tries, " tries")*/
			return tries
		}

	}
	return -1
}

func main() {
	fmt.Println(len(os.Args), os.Args)
	crypto := ""
	if len(os.Args) > 1 {
		crypto = os.Args[1]
	}
	useRSA := false
	if crypto == "rsa" {
		fmt.Println("Using RSA")
		useRSA = true
	} else if crypto == "eddsa" {
		fmt.Println("Using EDDSA")
		useRSA = false
	} else {
		fmt.Println("Please provide either \"rsa\" or \"eddsa\" as argument")
		return
	}

	targetCID := "bafkreidon73zkcrwdb5iafqtijxildoonbwnpv7dyd6ef3qdgads2jc4su"
	id_mh, _ := mh.FromB58String(targetCID)
	targetCIDKey := kspace.XORKeySpace.Key(id_mh)
	fmt.Printf("network_size, exp, tries, time[us]\n")
	exp_num := 100
	for size := 1000; size < 30001; size += 1000 {
		for i := 0; i < exp_num; i++ {
			start := time.Now()
			tries := generateNewKey(targetCIDKey, size, useRSA)
			elapsed := time.Since(start)
			fmt.Printf("%d,%d,%d,%d\n", size, i, tries, elapsed.Microseconds())
		}

	}

}
