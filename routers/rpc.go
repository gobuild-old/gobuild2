package routers

import (
	"net"
	"net/http"
	"net/rpc"

	"github.com/qiniu/log"
)

type GoRpc struct {
}

type Args struct {
	CgoEnabled bool
}

type Reply struct {
	OK     bool
	Repo   string
	Branch string
}

func (r *GoRpc) NewMission(args *Args, rep *Reply) error {
	log.Infof("cgo_enabled: %v", args.CgoEnabled)
	rep.Branch = "master"
	rep.Repo = "github.com/codeskyblue/fswatch"
	return nil
}

func GetMission(addr string, args *Args) (*Reply, error) {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil, err
	}
	reply := &Reply{}
	err = client.Call("GoRpc.NewMission", args, reply)
	return reply, err
}

func HandleRpc() {
	gr := new(GoRpc)
	rpc.Register(gr)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":8000")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
