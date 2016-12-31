package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	var (
		verbose bool
	)

	addr := flag.String("addr", ":8000", "http listen on")
	flag.BoolVar(&verbose, "V", false, "verbose")
	flag.Parse()

	ln := []byte("\n")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		if verbose {
			for k, v := range r.Header {
				log.Println("[DEBU]", k, v)
			}
		}

		log.Println("[DEBU]", r.Method, "BODY")
		fmt.Println("\033[35m" + string(content) + "\033[m")

		stdin := bufio.NewReader(os.Stdin)
		for {
			line, _, err := stdin.ReadLine()
			if err == nil {
				w.Write(append(line, ln...))
				continue
			}
			switch err {
			case io.EOF:
				return
			default:
				log.Println("[ERRO]", err)
				return
			}
		}
	})

	log.Println("[INFO] listen on", *addr)
	log.Println("[INFO] http server shutdown", http.ListenAndServe(*addr, nil))
}
