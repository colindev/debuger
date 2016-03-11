package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"golang.org/x/net/websocket"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "ws://host:port")
		os.Exit(1)
	}

	wsAddr := os.Args[1]

	conn, err := websocket.Dial(wsAddr, "", "http://localhost")
	if err != nil {
		panic(err)
	}

	quit := make(chan bool)
	go func() {
		var msg string
		for {
			select {
			case <-quit:
				break
			default:
				err := websocket.Message.Receive(conn, &msg)
				if err != nil {
					if err == io.EOF {
						// server close
						break
					}
					fmt.Println("Message receive error:", err)
					break
				}
				fmt.Println("<--", msg)
			}
		}
	}()

	cli := bufio.NewReader(os.Stdin)
	for {
		line, _, err := cli.ReadLine()
		if err != nil {
			fmt.Println(err)
			break
		}

		if err = websocket.Message.Send(conn, string(line)); err != nil {
			fmt.Println("Message send error:", err)
			break
		}
	}
	quit <- true
	os.Exit(0)

}
