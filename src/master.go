package main

import (
		"fmt"
		"net"
		"os"

		"actions"
)


// Entry point for server.
func main() {

	// Basic argument parsing. Accept a port and an optional log level.
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s port \n", os.Args[0])	
		os.Exit(1)
	}

	// Get the specified port.
	serverAddr, err := net.ResolveTCPAddr("tcp", ":" + os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", err)
		os.Exit(1)
	}

	// Start listening for incoming traffic from the ring.
	listenAddr, err := net.ResolveUDPAddr("udp", ":" + os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("My IP is: %v\n", actions.GetMyIP())

	// Launch all actions in separate goroutines.
	go actions.ListenForJoinRequests(serverAddr)
	go actions.MessagePrompt()
	go actions.ListenForRingMessage(listenAddr)

	// Spin until done.
	for {}
}
