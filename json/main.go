package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/colindev/gostdio"
)

func main() {

	err := gostdio.Read(os.Stdin, func(line []byte) error {
		buf := bytes.NewBuffer(nil)
		if err := json.Indent(buf, line, "", " "); err != nil {
			log.Println(err)
			return nil
		}

		fmt.Println(buf.String())

		return nil
	})

	if err != nil {
		log.Println(err)
	}

}
