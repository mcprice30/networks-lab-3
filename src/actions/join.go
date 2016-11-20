package actions

import (
	"fmt"
	"net"
	"os"

	"types"
)

// This thread will listen on TCP at the given port for any slaves to
// join the ring.
func ListenForJoinRequests(serverAddr *net.TCPAddr) {
	// Start listening on the given port.
	listener, err := net.ListenTCP("tcp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: %s\n", err)
		os.Exit(1)
	}

	nextHostAddr := GetMyIP()
	nextRID := byte(1)

	// When the server goes down, kill the connection.
	defer listener.Close()

	// Allocate buffer for holding requests.
	buf := make([]byte, 1024)

	for {
		// Get a request.
		conn, err := listener.AcceptTCP()
		fmt.Println()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		} else {
			fmt.Printf("Recieved TCP connection from %s\n", conn.RemoteAddr())
		}

		// Read and extract request bytes.
		if n, err := conn.Read(buf); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			conn.Close()
			continue
		} else if n != 3 {
			fmt.Fprintf(os.Stderr, "Invalid request size (%d)\n", n)
			conn.Close()
			continue
		} else {
			fmt.Printf("Recieved %s from client\n", showBytes(buf[:n]))
		}

		magicNumber := uint16(buf[1])<<8 + uint16(buf[2])
		if magicNumber != MagicNumber {
			fmt.Fprintf(os.Stderr, "Did not recieve '0x%04x', got 0x%04x\n",
				MagicNumber, magicNumber)
			conn.Close()
			continue
		}

		// Build the response.
		response := &types.ServerResponse{
			YourRID:     nextRID,
			NextSlaveIP: nextHostAddr,
		}

		respBytes, err := response.ToBytes()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not form response bytes: %s\n", err)
			conn.Close()
			continue
		}

		// Update local variables.
		remoteAddr := conn.RemoteAddr()
		fmt.Printf("REMOTE: %v\n", remoteAddr)
		if tcpAddr, ok := remoteAddr.(*net.TCPAddr); ok {
			nextHostAddr = tcpAddr.IP
			fmt.Println(nextHostAddr.String())
			nextPort := int(MY_GID)*5 + int(nextRID) + 10010
			nextAddr := fmt.Sprintf("%s:%d", nextHostAddr.String(), nextPort)
			outgoingAddr, err := net.ResolveUDPAddr("udp", nextAddr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not resolve next node: %s\n", err)
			}
			fmt.Printf("OUTGOING ADDR: %v\n", outgoingAddr)
			outgoingLock.Lock()
			outgoing, err = net.DialUDP("udp", nil, outgoingAddr)
			outgoingLock.Unlock()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Could not dial next node: %s\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Could not get client IP.\n")
		}
		nextRID++

		// Send the response to the client.
		fmt.Printf("Sending: %s\n", showBytes(respBytes))
		if n, err := conn.Write(respBytes); err != nil {
			fmt.Fprintf(os.Stderr, "Could not send response: %s\n", err)
		} else {
			fmt.Printf("Sent %d bytes to client.\n\n", n)
		}
		conn.Close()
		reprintChan <- true
	}
}
