package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

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

	quit := make(chan bool, 1)
	shutdown := make(chan os.Signal, 1)
	go func() {
		var msg string
		defer func() { shutdown <- syscall.SIGTERM }()
		for {
			select {
			case <-quit:
				quit <- true
				return
			default:
				err := websocket.Message.Receive(conn, &msg)
				if err != nil {
					if err == io.EOF {
						// server close
						fmt.Println("server closed")
						return
					}
					fmt.Println("Message receive error:", err)
					return
				}
				fmt.Println("<--", msg)
			}
		}
	}()

	go func() {
		cli := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-quit:
				break

			default:
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
		}
		quit <- true
	}()

	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	quit <- true
	os.Exit(0)

}
