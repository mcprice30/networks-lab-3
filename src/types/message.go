package types 

import ( "errors"
		"fmt"
)

const MagicNumber uint16 = 0x1234

type RingMessage struct {
	GID byte
	TTL byte
	RID_Dest byte
	RID_Src byte
	Msg string
	CheckSum byte
}

func MessageFromBytes(bytes []byte) (*RingMessage, error) {

	if computeCheckSum(bytes) != 0 {
		return nil, errors.New("Invalid checksum!")
	}

	out := &RingMessage{}		

	out.GID = bytes[0]
	if bytes[1] != byte((MagicNumber & 0xff) >> 8) || bytes[2] != byte(MagicNumber & 0xff) {
		return nil, errors.New("Magic Number does not match!")
	}

	out.TTL = bytes[3]
	out.RID_Dest = bytes[4]
	out.RID_Src = bytes[5]
	out.Msg = string(bytes[6:len(bytes)-1])
	// don't compute checksum, it will be handled when we spit out the bytes.
	return out, nil
}

func computeCheckSum(bytes []byte) byte {
	sum := uint16(0)
	for _, b := range bytes {
		sum += uint16(b)
		sum = (sum & 0xff00) >> 8 + (sum & 0x00ff)
	}
	sum = (sum & 0xff00) >> 8 + (sum & 0x00ff)
	return byte(sum & 0xff) ^ 0xff
}

func (rm *RingMessage) ToBytes() ([]byte) {
	out := make([]byte, 7 + len(rm.Msg))
	fmt.Println(len(out))
	out[0] = rm.GID
	out[1] = byte((MagicNumber & 0xff00) >> 8)
	out[2] = byte(MagicNumber & 0xff)
	out[3] = rm.TTL
	out[4] = rm.RID_Dest
	out[5] = rm.RID_Src
	for i, c := range rm.Msg {
		out[i + 6] = byte(c)	
	}
	fmt.Println(out)
	checkSum := computeCheckSum(out)
	out[len(out)-1] = checkSum
	//fmt.Println(computeCheckSum(out))
	return out
}

func (rm *RingMessage) String() string {
	return fmt.Sprintf("Sending '%s' from %d to %d with TTL %d", rm.Msg,
		rm.RID_Src, rm.RID_Dest, rm.TTL)
}
