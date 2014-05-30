// handle gobuild.io hooks
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/log"
	"io"
	"io/ioutil"
	"net/http"
)

func Hello(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello")
}

const (
	AUTHCHECKER = "./authchecker"
	RECEIVER    = "./receiver"
)

func Hook(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	info := new(HookInfo)
	err := json.Unmarshal(body, info) // body)
	if err != nil {
		log.Errorf("unmarshal request hook infomation error: %v", err)
		return
	}
	if sh.Test("x", AUTHCHECKER) {
		err := sh.Command(AUTHCHECKER, info.Secret).Run()
		if err != nil {
			log.Errorf("authchecker not passed: %v", err)
			return
		}
	}
	if !sh.Test("x", RECEIVER) {
		log.Errorf("need script '%s'", RECEIVER)
		return
	}
	if err := sh.Command(RECEIVER, info.Repo, info.ZipballUrl).Run(); err != nil {
		log.Errorf("call %s error: %v", RECEIVER, err)
	}
	fmt.Printf("receive: %v\n", string(body))
}

func init() {
	http.HandleFunc("/", Hello)
	http.HandleFunc("/webhook", Hook)
}

func main() {
	addr := flag.String("http", ":8877", "HTTP service address")
	hookby := flag.String("hookby", "gogs", "hook by which service")
	flag.Parse()

	log.Printf("Hook for %s, Listening on %s\n", *hookby, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
