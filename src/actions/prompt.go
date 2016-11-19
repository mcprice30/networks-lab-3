package actions

import (
		"bufio"
		"fmt"
		"os"
		"strconv"
		"strings"

		"types"
)

var reprintChan chan bool = make(chan bool)

func ReprintPrompt() {
	for {
		<- reprintChan	
		fmt.Print("Enter message: ")
	}
}

func MessagePrompt() {
	go ReprintPrompt()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		fmt.Print("Enter target RID: ")
		ridNum, _ := reader.ReadString('\n')
		targetRid, _ := strconv.Atoi(ridNum)

		message := &types.RingMessage {
			GID: MY_GID,
			TTL: 10,
			RID_Dest: byte(targetRid),
			RID_Src: 0,
			Msg: text,
		}

		fmt.Println(message.Msg)
	
		sendMessage(message)
	}
}
