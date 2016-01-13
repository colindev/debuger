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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		log.Debug(r.Method, "BODY", string(content))

		line, _, err := bufio.NewReader(os.Stdin).ReadLine()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(line)
	})

	log.Info("listen on", *addr)
	log.Info(http.ListenAndServe(*addr, nil))
}
