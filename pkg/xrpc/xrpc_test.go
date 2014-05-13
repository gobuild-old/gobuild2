package xrpc

import (
	"log"
	"net/http"
	"testing"
	"time"
)

func init() {
	HandleRpc()
	go http.ListenAndServe(":12345", nil)
	time.Sleep(time.Millisecond * 10)
	DefaultServer = "localhost:12345"
}

func TestRpcCall(t *testing.T) {
	args := new(Args)
	args.Arch = "386"
	reply, err := GetMission(args)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)

	err = UpdateStatus(args)
	if err != nil {
		log.Fatal(err)
	}
}
