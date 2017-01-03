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
	"sync"
)

func main() {

	var (
		verbose bool
		addr    string
		lc      sync.Mutex
	)

	flag.StringVar(&addr, "addr", ":8000", "http listen on")
	flag.BoolVar(&verbose, "V", false, "verbose")
	flag.Parse()

	ln := []byte{'\r', '\n'}
	// 清空之前的輸入
	stdin := bufio.NewReader(os.Stdin)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lc.Lock()
		defer lc.Unlock()
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Fatal("response writer can't assert to http.Flusher")
		}

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		fmt.Printf("%s %s\n", r.Method, r.URL)
		if verbose {
			for k, v := range r.Header {
				fmt.Printf("\033[2;33m%s:\033[m %v\n", k, v)
			}
		}

		fmt.Println("\033[2;35m--- body start\n" + string(content) + "\n--- body end\033[m")
		fmt.Println("--- input response body")
		defer fmt.Println("--- finish response")

		for {
			line, _, err := stdin.ReadLine()
			if err == nil {
				w.Write(append(line, ln...))
				flusher.Flush()
				if len(line) == 0 {
					return
				}
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

	log.Println("[INFO] listen on", addr)
	log.Println("[INFO] http server shutdown", http.ListenAndServe(addr, nil))
}
