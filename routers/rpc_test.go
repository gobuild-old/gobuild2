package routers

import (
	"net/http"
	"testing"
	"time"
)

func init() {
	HandleRpc()
	go http.ListenAndServe(":12345", nil)
	time.Sleep(time.Millisecond * 10)
}

func TestRpcCall(t *testing.T) {
	args := new(Args)
	args.Arch = "386"
	reply, err := GetMission("localhost:12345", args)
	if err != nil {
		t.Error(err)
	}
	t.Log(reply)
}
