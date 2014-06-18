// handle gobuild.io hooks
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/log"
)

func Hello(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "hello")
}

const (
	AUTHCHECKER = "./authchecker"
	RECEIVER    = "./receiver"
)

func HookGogs(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	payload := new(GogsPayload)
	err := json.Unmarshal(body, payload) // body)
	if err != nil {
		log.Errorf("unmarshal request hook infomation error: %v", err)
		return
	}
	if sh.Test("x", AUTHCHECKER) {
		err := sh.Command(AUTHCHECKER, payload.Secret).Run()
		if err != nil {
			log.Errorf("authchecker not passed: %v", err)
			return
		}
	}
	if !sh.Test("x", RECEIVER) {
		log.Errorf("need script '%s'", RECEIVER)
		return
	}
	env := map[string]string{
		"PUSH_NAME":  payload.Pusher.Name,
		"PUSH_EMAIL": payload.Pusher.Email,
	}
	if err := sh.Command(RECEIVER, payload.Ref, env).Run(); err != nil {
		log.Errorf("call %s error: %v", RECEIVER, err)
	}
	fmt.Printf("receive: %v\n", payload.Ref)
}

func init() {
	http.HandleFunc("/", Hello)
	//http.HandleFunc("/webhook", HookGogs)
	http.HandleFunc("/webhook/gogs", HookGogs)
}

func main() {
	addr := flag.String("http", ":8877", "HTTP service address")
	hookby := flag.String("hookby", "gogs", "hook by which service")
	flag.Parse()

	log.Printf("Hook for %s, Listening on %s\n", *hookby, *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
