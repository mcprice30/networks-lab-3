package actions

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"types"
)

// reprintChan listens for any requests to reprint the prompt.
var reprintChan chan bool = make(chan bool)

// ReprintPrompt will reprint the beginning of the prompt, in the event that
// another message is sent in the mean time.
func ReprintPrompt() {
	for {
		<-reprintChan
		fmt.Print("Enter message: ")
	}
}

// MessagePrompt will prompt a user for a message and an RID. Upon recieving
// one, it will send the corresponding message in the ring.
func MessagePrompt() {
	go ReprintPrompt()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		fmt.Print("Enter target RID: ")
		ridNum, _ := reader.ReadString('\n')
		targetRid, err := strconv.Atoi(strings.TrimSpace(ridNum))
		for err != nil {
			fmt.Println(err)	
			fmt.Print("Enter target RID: ")
			ridNum, _ := reader.ReadString('\n')
			targetRid, err = strconv.Atoi(strings.TrimSpace(ridNum))
		}

		message := &types.RingMessage{
			GID:      MY_GID,
			TTL:      5,
			RID_Dest: byte(targetRid),
			RID_Src:  0,
			Msg:      text,
		}
		sendMessage(message)
	}
}
