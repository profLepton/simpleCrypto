package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Initializing global variables
var (
	DIFFICULTY     = 6
	FILE_NAME      = "./blockChain.txt"
	GENESIS_BLOCK  = "0000000000000000000000000000000000000000000000000000000000000000~satoshi~11970128322"
	TAIL           = GENESIS_BLOCK
	TERMINATOR     = '\n'
	ADDRESS        = ":8888"
	BLOCK_REQUEST  = "TAIL"
	BLOCK_ACCEPTED = "BLOCK ACCEPTED"
)

func main() {

	log.Println("Server> Starting the server.")

	//Check and create BlockChain File
	TAIL = initializeServer(FILE_NAME)

	//Start server and handle connections
	log.Println("Server> Listening for connections.")
	serverListener, err := net.Listen("tcp", ADDRESS)
	if err != nil {
		panic(err)
	}

	//Run and look for connections.
	for {
		serverConnection, err := serverListener.Accept()
		if err != nil {
			panic(err)
		}

		go HandleConnections(serverConnection)
	}

}

// Block is defined as the {hash of previous block}~{hash of name}~{name}~{nonce}
// Store the entire block chain into a file
// Each block has the hash of the previous b
// When the server receives a block, it checks if the hash of the block string
// \n satisfies the given difficulty rating.
// The difficulty rating is constant
// The server needs to be able to connect to multiple nodes asynchronously,
// It needs to be able to verify 2 things, namely the 1) If the block string contains the
// hash of the current blockchain tail, second if the hash of the block string satisfies the POW requirement.
// Two functions needed HashString, veriifyPOW
// Need to story the current block tip
// Also need to give entire block chain to all  nodes when requested.
// Need read write permissions to file.
// Need to update the file and blockchain tip when new block is accepted.
// Need to convert hash to hexa

// Returns the Hex of a Hash as a string.
func HexaHashFromString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%064x", hash)
}

// Handles connections from all clients
// Recognize different requests - Block requests & Submit requests
// Be able to read and answer block requests
// Be able to verify and save submit requests,
func HandleConnections(conn net.Conn) {
	log.Printf("Server> Connected to %s\n", conn.RemoteAddr().String())

	//Read incoming connections
	bufReader := bufio.NewReader(conn)

	recvdBytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		log.Printf("Server> TCP error: %s\n", err.Error())
	}
	recvdString := string(recvdBytes)
	var sendBytes []byte
	var sendString string
	if strings.HasPrefix(recvdString, BLOCK_REQUEST) {
		log.Print("Server> Received a block request.")
		sendString = TAIL
	} else {
		newBl := recvdString[:len(recvdString)-1]
		log.Print("Server> Received block\n")
		if VerifyBlockString(TAIL, newBl) {
			//Block has been accepted.
			log.Print("Server> New block has been accepted!")
			TAIL = newBl
			appendBlockToDisk(newBl, FILE_NAME)
			sendString = "BLOCK ACCEPTED"
		} else {
			log.Print("Server> Invalid submission.")
			sendString = "Invalid Submission"
		}
	}

	//Sending stuff
	sendBytes = []byte(fmt.Sprintf("%s\n", sendString))
	_, err = conn.Write(sendBytes)
	if err != nil {
		log.Printf("Server> TCP error: %s\n", err.Error())
	}

	//Disconnect
	conn.Close()
}

// Takes in a string and a difficulty int and checks if the POW is satisfied.
func checkPOW(s string, difficulty int) bool {
	zeroes := 0
	for i := range make([]int, difficulty) {
		if s[i] == '0' {
			zeroes++
		}
	}
	return zeroes == difficulty
}

// Verifies if the submitted block is a valid solution.
func VerifyBlockString(TAIL string, submittedBlock string) bool {
	bctHash := HexaHashFromString(TAIL)
	block_array := strings.Split(submittedBlock, "~")

	//Check if block has this bctHash at the correct spot, if yes then check pow

	//else return false if pow is not satisfied or hctHash is misplaced or incorrect.
	for i := range block_array {
		log.Print(block_array[i])
	}
	if bctHash == block_array[0] {
		//check proof of work
		submittedBlockHash := HexaHashFromString(submittedBlock)
		return checkPOW(submittedBlockHash, DIFFICULTY)

	}
	return false
}

func initializeServer(file_name string) string {
	f, err := os.Open(file_name)
	if err != nil {
		f, err = os.Create(file_name)
		if err != nil {
			panic(err)
		}
		f.WriteString(TAIL + "\n")
		return TAIL
	}

	scanner := bufio.NewScanner(f)

	var block string

	for scanner.Scan() {
		block = scanner.Text()
	}
	f.Close()
	return block
}

func appendBlockToDisk(block string, file_name string) {
	f, err := os.OpenFile(FILE_NAME, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	f.WriteString(block + "\n")
	log.Print("Server> Updating block chain")
	f.Close()
}
