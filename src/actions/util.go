package actions

import (
		"fmt"
		"net"
		"os"
		"sync"

		"types"
)

const MY_GID byte = 11
const MagicNumber uint16 = 0x1234

var outgoingLock *sync.Mutex = &sync.Mutex{}
var outgoing *net.UDPConn

func sendMessage(datagram *types.RingMessage) {
	msgBytes := datagram.ToBytes()
	outgoingLock.Lock()
	defer outgoingLock.Unlock()
	if outgoing == nil {
		fmt.Fprintln(os.Stderr, "No node to send to!")
		return
	}

	if _, err := outgoing.Write(msgBytes); err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: ", err)
		return
	}
	fmt.Printf("Sent bytes: %v", msgBytes)
}

func GetMyIP() net.IP {
	ifaces, err := net.Interfaces()
	if err != nil { return nil }
	var out net.IP

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if (out == nil || out.IsLoopback()) {
					out = v.IP
				}
			case *net.IPAddr:
				if (out == nil || out.IsLoopback()) {
					out = v.IP
				}
			}
		}
	}
	return out
}

func showBytes(in []byte) string {
	out := "["
	for i, b := range in {
		if i > 0 { out += " " }
		out += fmt.Sprintf("%02x", b)
	}
	return (out + "]")
}
