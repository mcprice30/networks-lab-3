package actions

import (
	"fmt"
	"net"
	"os"

	"types"
)

// ListenForRingMessage will listen for messages incoming from the token ring.
// It will then forward the message to the appropriate node.
func ListenForRingMessage(addr *net.UDPAddr) {
	incoming, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot listen for ring messages (%s)\n", err)
		return
	}

	for {
		buffer := make([]byte, 1024)
		var n int
		var err error

		fmt.Println() // Get some separation from the prompt.
		if n, err = incoming.Read(buffer); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read ring messages (%s)\n", err)
		} else {
			fmt.Printf("GOT %s\n", showBytes(buffer[:n]))
		}

		// Build the message. This will handle error detection.
		message, err := types.MessageFromBytes(buffer[:n])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error in message recieved (%s)\n", err)
			continue
		}

		fmt.Println("Recieved: ", message)

		// If it is for us, print it out.
		if message.RID_Dest == 0 {
			fmt.Printf("Got message - '%s' from client with RID %d and GID %d\n",
				message.Msg, message.RID_Src, message.GID)
		} else { // Otherwise, forward it.
			message.TTL--
			if message.TTL > 0 {
				sendMessage(message)
			} else {
				fmt.Println("Dropping message, TTL too large")
			}
		}
		reprintChan <- true // Reprint the prompt.
	}
}
