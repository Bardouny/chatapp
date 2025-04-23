package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Message struct {
	To      byte   `json:"To"`
	Message string `json:"Message"`
}

func main() {
	conn, err := net.Dial("tcp", ":5000")

	if err != nil {
		log.Fatal("error at client")
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	buffer := make([]byte, 1024)
	go func() {
		for {
			n, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("error")
				log.Fatal("Server is over")
			}

			fmt.Println("server : ", string(buffer[:n]))
		}

	}()

	for {
		var msg = Message{}
		fmt.Println("Sending to : ")
		to, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("error")
		}
		fmt.Println("Message : ")
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("error")
		}

		time.Sleep(time.Second)

		msg.To = to[0]
		msg.Message = data

		res, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("faild to send message")
		}

		conn.Write(res)

	}

}
