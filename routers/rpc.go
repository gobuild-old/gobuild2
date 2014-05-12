package routers

import (
	"net/rpc"

	"github.com/qiniu/log"
)

type Rpc struct {
}

type Args struct {
	Os, Arch string
	Host     string
}

type Reply struct {
	OK     bool
	Repo   string
	Branch string
}

func (r *Rpc) NewMission(args *Args, rep *Reply) error {
	log.Infof("arch: %v", args.Arch)
	log.Infof("host: %v", args.Host)
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
	err = client.Call("Rpc.NewMission", args, reply)
	return reply, err
}

func HandleRpc() {
	gr := new(Rpc)
	rpc.Register(gr)
	rpc.HandleHTTP()
}
