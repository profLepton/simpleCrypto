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

var (
	DIFFICULTY     = 3
	HEADER_FILE    = "./BlockHeader.txt"
	BLOCK_FILE     = "./BlockChain.txt"
	FRIENDS        = "FRIENDS.txt"
	GENESIS_BLOCK  = "0000000000000000000000000000000000000000000000000000000000000000~satoshi~11970128322"
	TAIL           = GENESIS_BLOCK
	TERMINATOR     = '\n'
	ADDRESS        = ":8888"
	BLOCK_REQUEST  = "TAIL"
	BLOCK_ACCEPTED = "BLOCK ACCEPTED"
	DELIMITTER     = "~"
)

func main() {
	log.Println("Node> Starting node.")
	var TAIL = InitializeServer(HEADER_FILE)
	log.Printf("Node> Tail is %s\n", TAIL)

	//Start server and handle connections
	log.Println("Node> Listening for connections.")
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

//Difficulty update function

//Block header contains time of block
//Block header contains merkle tree of transactions
//Block is of len 1024 bytes.

// Creates/Opens the block header file and ip address file and returns the block tip.
func InitializeServer(HEADER_FILE string) string {

	f, err := os.Open(HEADER_FILE)
	if err != nil {
		f, err = os.Create(HEADER_FILE)
		if err != nil {
			panic(err)
		}
		f.WriteString(GENESIS_BLOCK + "\n")
		return GENESIS_BLOCK
	}

	scanner := bufio.NewScanner(f)

	var block string

	for scanner.Scan() {
		block = scanner.Text()
	}
	f.Close()

	_, err = os.Open(FRIENDS)
	if err != nil {
		_, err = os.Create(FRIENDS)
		if err != nil {
			panic(err)
		}
	}
	return block
}

func HandleConnections(conn net.Conn) {
	log.Printf("Node> Received connection from %s\n", conn.RemoteAddr().String())

	//Read incoming connections
	bufReader := bufio.NewReader(conn)

	recvdBytes, err := bufReader.ReadBytes('\n')
	if err != nil {
		log.Printf("Node> TCP error: %s\n", err.Error())
	}
	recvdString := string(recvdBytes)
	var sendBytes []byte
	var sendString string
	if strings.HasPrefix(recvdString, BLOCK_REQUEST) {
		log.Print("Node> Received a block request.")
		sendString = TAIL
	} else {
		newBl := recvdString[:len(recvdString)-1]
		log.Print("Server> Received block\n")
		if VerifyBlockString(TAIL, newBl) {
			//Block has been accepted.
			log.Print("Server> New block has been accepted!")
			TAIL = newBl
			AppendHeaderToDisk(newBl, HEADER_FILE)
			sendString = "BLOCK ACCEPTED"
		} else {
			log.Print("Server> Invalid submission.")
			sendString = "INVALID SUBMISSION"
		}
	}

	//Sending stuff
	sendBytes = []byte(fmt.Sprintf("%s%c", sendString, TERMINATOR))
	_, err = conn.Write(sendBytes)
	if err != nil {
		log.Printf("Node> TCP error: %s\n", err.Error())
	}

	//Disconnect
	conn.Close()
}

func AppendHeaderToDisk(block string, file_name string) {
	f, err := os.OpenFile(HEADER_FILE, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	f.WriteString(block + "\n")
	log.Print("Server> Updating block chain")
	f.Close()
}

func VerifyBlockString(TAIL string, submittedBlock string) bool {
	bctHash := HexaHashFromString(TAIL)
	block_array := strings.Split(submittedBlock, DELIMITTER)

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

func checkPOW(s string, difficulty int) bool {
	zeroes := 0
	for i := range make([]int, difficulty) {
		if s[i] == '0' {
			zeroes++
		}
	}
	return zeroes == difficulty
}

func HexaHashFromString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%064x", hash)
}
