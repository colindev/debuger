package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	var (
		ip         string
		port       string
		listenMode bool
		verbose    bool
		network    = "tcp"
	)
	cli := flag.CommandLine
	cli.BoolVar(&listenMode, "l", false, "listen mode")
	cli.BoolVar(&verbose, "v", false, "verbose")
	cli.Usage = func() {
		fmt.Printf("Usage: command [options] ip port\n")
		cli.PrintDefaults()
	}
	cli.Parse(os.Args[1:])
	args := cli.Args()
	if len(args) != 2 {
		cli.Usage()
		os.Exit(2)
	}

	if verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	ip, port = args[0], args[1]

	conn, err := getConn(ip, port, listenMode, network)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		w := bufio.NewWriter(conn)
		bufStdin := bufio.NewReader(os.Stdin)
		for {
			line, _, err := bufStdin.ReadLine()
			if err != nil {
				log.Fatal(err)
			}

			w.Write(line)
			w.WriteByte('\n')
			if err := w.Flush(); err != nil {
				log.Println(err)
			}
		}
	}()

	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(line))
	}

}

func getConn(ip, port string, listenMode bool, network string) (net.Conn, error) {

	if !listenMode {
		return net.Dial(network, fmt.Sprintf("%s:%s", ip, port))
	}

	switch network {
	case "tcp":
		addr, err := net.ResolveTCPAddr(network, fmt.Sprintf("%s:%s", ip, port))
		if err != nil {
			return nil, err
		}

		listener, err := net.ListenTCP(network, addr)
		if err != nil {
			return nil, err
		}
		defer listener.Close()

		return listener.Accept()

		//	case "udp":
		//		addr, err := net.ResolveUDPAddr(network, fmt.Sprintf("%s:%s", ip, port))
		//		if err != nil {
		//			return nil, err
		//		}
		//
		//		return net.ListenUDP(network, addr)
	}

	return nil, fmt.Errorf("unkown network")
}
