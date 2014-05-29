// handle gobuild.io hooks
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func Hello(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello")
}

func Hook(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("receive: %v\n", string(body))
}

func init() {
	http.HandleFunc("/", Hello)
	http.HandleFunc("/hook", Hook)
}

func main() {
	addr := flag.String("http", ":8877", "HTTP service address")
	flag.Parse()

	log.Printf("Listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
