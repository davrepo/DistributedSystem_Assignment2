package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	proto "simpleGuide/grpc"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id int
}

var (
	serverPorts = flag.String("sPorts", "", "Comma-separated server port numbers")
)

func main() {
	// Configure log output to file
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)
	defer logFile.Close()

	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id: 1,
	}

	// Wait for the client (user) to ask for the time
	go waitForTimeRequest(client)

	// Create a channel to block the main thread
	done := make(chan bool)

	// Block the main thread
	<-done
}

func waitForTimeRequest(client *Client) {
	// Connect to multiple servers
	serverConnections := connectToMultipleServers()

	var serverTimes []time.Time

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Iterate through each server connection
		for port, conn := range serverConnections {
			// Measure roundtrip time
			startTime := time.Now()

			// Ask the server for the time
			timeReturnMessage, err := conn.AskForTime(context.Background(), &proto.AskForTimeMessage{
				ClientId: int64(client.id),
			})

			// Calculate roundtrip time
			roundtripTime := time.Since(startTime)

			if err != nil {
				log.Printf("Error from server at port %d: %s", port, err.Error())
			} else {
				// Implement Christian's Algorithm to adjust the server time
				serverTime, _ := time.Parse(time.RFC3339, timeReturnMessage.Time)
				adjustedServerTime := serverTime.Add(roundtripTime / 2)
				serverTimes = append(serverTimes, adjustedServerTime)

				log.Printf("Server at port %d says the adjusted time is %s. Roundtrip Time: %d ms", port, adjustedServerTime, roundtripTime.Milliseconds())
			}
		}

		// Calculate the offset between server times
		if len(serverTimes) >= 2 {
			timeDiff := serverTimes[1].Sub(serverTimes[0])
			log.Printf("Time difference (offset) between server 1 and server 2: %v", timeDiff)
		}
	}
}

func connectToMultipleServers() map[int]proto.TimeAskClient {
	serverPortsArray := strings.Split(*serverPorts, ",")
	serverConnections := make(map[int]proto.TimeAskClient)

	for _, portStr := range serverPortsArray {
		port, _ := strconv.Atoi(portStr)
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(port), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("Could not connect to port %d: %s", port, err.Error())
		} else {
			log.Printf("Connected to the server at port %d", port)
			serverConnections[port] = proto.NewTimeAskClient(conn)
		}
	}
	return serverConnections
}
