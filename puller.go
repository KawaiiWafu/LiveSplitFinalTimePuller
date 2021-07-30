package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func connectionErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Setup name
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	runner, _ := reader.ReadString('\n')

	// Setup API server connection
	started := false
	fmt.Print("Enter API server address: ")
	server, _ := reader.ReadString('\n')
	server = strings.TrimSuffix(server, "\r\n")

	// LiveSplit Server connection
	host := "localhost:16834"
	connection, err := net.Dial("tcp", host)
	connectionErr(err)

	// Messages
	phaseMsg := []byte("getcurrenttimerphase\r\n")
	endTimeMessage := []byte("getfinaltime\r\n")
	end := "Ended\r\n"
	running := "Running\r\n"

	for {
		// Reading for "Ended" event
		connection.Write(phaseMsg)
		reply := make([]byte, 4096)
		connection.Read(reply)
		reply = bytes.Trim(reply, "\x00")

		// Start timer via API
		if !started && string(reply) == running {
			_, err := http.Get("http://" + server + "/nodecg-mafiamarathon/startTimer")
			connectionErr(err)
			started = true
		}

		// Send to API on "Ended" event
		if string(reply) == end {
			// Read final time
			connection.Write(endTimeMessage)
			reply := make([]byte, 4096)
			connection.Read(reply)
			reply = bytes.Trim(reply, "\x00")
			finalTime := strings.TrimSuffix(string(reply), "\r\n")

			// Send API request
			params := url.Values{}
			params.Add("runner", strings.TrimSuffix(runner, "\r\n"))
			params.Add("time", finalTime)
			_, err := http.Get("http://" + server + "/nodecg-mafiamarathon/endRun?" + params.Encode())
			connectionErr(err)
			break
		}
	}
}
