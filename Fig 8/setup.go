package main

import (
	"encoding/csv"
	"log"
	"os"
	"os/exec"
	"strings"
	"fmt"
)

// add the file to IPFS from the node - ensure ipfs daemon is running
func getCID(filename string) []byte {
	out, err := exec.Command("ipfs", "add", filename).Output()
	if err != nil {
		log.Println("getCID() failed!")
		log.Println(err)
		return nil
	}
	return out
}



func main(){

  // specify the number of files to be created, ensure empty file 'initial.txt' exists
	cidnumber := 100
	fmt.Println("Generating", cidnumber, "files...")
	fmt.Println()
	initialFile := "./data/initial.txt"
  
	// get CID and add to the network
	initialCID := string(getCID(initialFile))
	initialCID = strings.ReplaceAll(initialCID, "\n", "")
	initialCID2 := strings.Split(initialCID, " ")
	fmt.Printf("First CID: '%s'\n", initialCID2[1])

	cidList, err := os.OpenFile("./data/cidList.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
		fmt.Println("error")
	}
	defer cidList.Close()

	cidListwriter := csv.NewWriter(cidList)
	record := []string{initialCID2[1]}
	cidListwriter.Write(record)

  // add new cids to the network and log in 'cidList.csv'
	for i:=0;i<cidnumber;i++{
		f, err := os.OpenFile("./data/initial.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
			fmt.Println("error")

		}
		defer f.Close()
		if _, err := f.WriteString("foo bar //d 2443 \n"); err != nil {
			log.Println(err)
			fmt.Println("error")

		}
		newCid := string(getCID(initialFile))
		newCid = strings.ReplaceAll(newCid, "\n", "")
		newCid2 := strings.Split(newCid, " ")
		fmt.Printf("New CID: '%s'\n", newCid2[1])
		
		record := []string{newCid2[1]}
		cidListwriter.Write(record)
	}

	cidListwriter.Flush()

}
