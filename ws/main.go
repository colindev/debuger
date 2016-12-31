package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

// Headers of request
type Headers map[string][]string

// Set header
func (h Headers) Set(v string) error {
	s := strings.SplitN(v, ":", 2)
	if len(s) != 2 {
		return errors.New("key:value")
	}
	h[s[0]] = append(h[s[0]], s[1])
	return nil
}

func (h Headers) String() string {
	return fmt.Sprintf("%+v", (map[string][]string)(h))
}

func main() {

	var (
		listenMode bool
		headers    Headers
		verbose    bool
	)

	cli := flag.CommandLine
	cli.Usage = func() {
		fmt.Println("Usage: command [option] ws[s]://host[:port][/path]")
		cli.PrintDefaults()
	}
	cli.BoolVar(&verbose, "v", false, "verbose")
	cli.BoolVar(&listenMode, "l", false, "listen mode")
	cli.Var(&headers, "H", "header")
	cli.Parse(os.Args[1:])
	args := cli.Args()
	if len(args) != 1 {
		cli.Usage()
		os.Exit(2)
	}

	conn, err := getConn(args[0], listenMode, headers)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(p))
		}
	}()

	bufStdin := bufio.NewReader(os.Stdin)
	for {
		line, _, err := bufStdin.ReadLine()
		if err != nil {
			log.Fatal(err)
		}

		if err := conn.WriteMessage(websocket.TextMessage, line); err != nil {
			log.Fatal(err)
		}
	}
}

func getConn(addr string, listenMode bool, headers Headers) (conn *websocket.Conn, err error) {

	if !listenMode {
		conn, _, err = websocket.DefaultDialer.Dial(addr, http.Header(headers))
		return
	}

	s := strings.SplitN(addr, "://", 2)
	listener, err := net.Listen("tcp", s[len(s)-1])
	if err != nil {
		log.Fatal(err)
	}
	upgrader := websocket.Upgrader{}
	http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err = upgrader.Upgrade(w, r, nil)
		listener.Close()
	}))

	return
}
