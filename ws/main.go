package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/websocket"
)

func main() {

	var (
		wsAddr string
		origin string
	)

	flag.StringVar(&origin, "origin", "http://localhost", "client origin")
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Println("Usage:", os.Args[0], "ws[s]://host[:port]")
		os.Exit(1)
	}
	wsAddr = flag.Arg(0)

	conn, err := websocket.Dial(wsAddr, "", origin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
				goto quit

			default:
				line, _, err := cli.ReadLine()
				if err != nil {
					fmt.Println(err)
					goto quit
				}

				if err = websocket.Message.Send(conn, string(line)); err != nil {
					fmt.Println("Message send error:", err)
					goto quit
				}
			}
		}
	quit:
		quit <- true
	}()

	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	quit <- true
	os.Exit(0)

}
