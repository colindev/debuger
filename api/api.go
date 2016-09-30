package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	addr := flag.String("addr", ":8000", "http listen on")
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	c := make(chan []byte)

	go func() {

		r := bufio.NewReader(os.Stdin)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				log.Println("[ERRO]", err)
				close(c)
				return
			}
			c <- line
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if *verbose {
			for k, v := range r.Header {
				log.Println("[DEBU]", k, v)
			}
		}

		log.Println("[DEBU]", r.Method, "BODY", string(content))

		line := <-c

		log.Println("[DEBU] READ STDIN", string(line))
		w.Write(line)
	})

	log.Println("[INFO] listen on", *addr)
	log.Println("[INFO] http server shutdown", http.ListenAndServe(*addr, nil))
}
