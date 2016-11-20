package types

import (
	"fmt"
	"net"
)

const MY_GID byte = 11

// ServerRespose is the message returned from the master whenever a slave joins
// the ring.
type ServerResponse struct {
	YourRID     byte
	NextSlaveIP net.IP
}

// Produce a byte slice from the response.
func (sR *ServerResponse) ToBytes() ([]byte, error) {
	out := make([]byte, 8)
	out[0] = MY_GID
	out[1] = byte((MagicNumber & 0xff00) >> 8)
	out[2] = byte(MagicNumber & 0xff)
	out[3] = sR.YourRID

	if ip := sR.NextSlaveIP.To4(); len(ip) != 4 {
		return nil, fmt.Errorf("Expected IPv4 address, got: %s", ip)
	} else {
		for i, b := range ip {
			out[4+i] = b
		}
	}
	return out, nil
}
