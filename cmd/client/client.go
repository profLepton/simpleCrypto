package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"strconv"
)

var (
	DIFFICULTY     = 6
	ADDRESS        = ":8888"
	TERMINATOR     = []byte("\n")
	BLOCK_REQUEST  = "TAIL\n"
	BLOCK_ACCEPTED = "BLOCK ACCEPTED"
	PAYLOAD        = "Akshay's First Coin"
)

func main() {

	log.Print("Client> Starting client.\n")
	//Connect to server, get tail of blockChain, close conn
	//Create a block, update nonce until difficulty is satisfied.
	//Connect to server, send the block and SUCCESS!

	conn, err := net.Dial("tcp", ADDRESS)
	if err != nil {
		panic(err)
	}

	//Request the tip
	_, err = conn.Write([]byte(BLOCK_REQUEST))
	if err != nil {
		panic(err)
	}
	//Read the tip
	bufReader := bufio.NewReader(conn)
	rcvdBytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		panic(err)
	}
	conn.Close()

	rcvdString := string(rcvdBytes)
	rcvdString = rcvdString[:len(rcvdString)-1]
	log.Printf("Client> Received block: %s\n", rcvdString)
	tail_hash := HexaHashFromString(rcvdString)

	s := mine(tail_hash, PAYLOAD, DIFFICULTY)

	sendString := s + "\n"

	log.Print(sendString)
	conn, err = net.Dial("tcp", ADDRESS)
	if err != nil {
		panic(err)
	}
	//Submit Block
	_, err = conn.Write([]byte(sendString))
	if err != nil {
		panic(err)
	}
	//Read response
	bufReader = bufio.NewReader(conn)
	rcvdBytes, err = bufReader.ReadBytes('\n')
	conn.Close()
	rcvdString = string(rcvdBytes)
	log.Print(rcvdString)

}

// Returns the Hex of a Hash as a string.
func HexaHashFromString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%064x", hash)
}

// Checks if the string has the required number of zeroes.
func checkPOW(s string, difficulty int) bool {
	zeroes := 0
	for i := range make([]int, difficulty) {
		if s[i] == '0' {
			zeroes++
		}
	}
	return zeroes == difficulty
}

func mine(tailHash string, payload string, difficulty int) string {
	i := 0
	run := true
	var trial_block string
	var trial_block_hash string
	for run {
		trial_block = tailHash + "~" + payload + "~" + strconv.Itoa(i)
		trial_block_hash = HexaHashFromString(trial_block)
		run = !checkPOW(trial_block_hash, difficulty)
		i++
	}
	log.Printf("Client> Found good nonce: %d \n The Hash is %s \n", i, trial_block_hash)
	return trial_block
}
