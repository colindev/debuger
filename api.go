package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ngaut/log"
)

func main() {

	addr := flag.String("addr", ":8000", "http listen on")
	flag.Parse()

	c := make(chan []byte)

	go func() {

		r := bufio.NewReader(os.Stdin)
		for {
			line, _, err := r.ReadLine()
			if err != nil {
				log.Debug(err)
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

		log.Debug(r.Method, "BODY", string(content))

		line := <-c
		log.Debug("READ STDIN", string(line))
		w.Write(line)
	})

	log.Info("listen on", *addr)
	log.Info(http.ListenAndServe(*addr, nil))
}
