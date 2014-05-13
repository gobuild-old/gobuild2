package xrpc

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/qiniu/log"
)

var DefaultServer = "localhost:8000"

type Rpc struct{}

const (
	ST_PENDING    = "pending"
	ST_BUILDING   = "building"
	ST_PUBLISHING = "publishing"
	ST_DONE       = "done"
	ST_ERROR      = "error"
)

type Args struct {
	Os, Arch string
	Host     string
	Mid      int64 // mission id
	Status   string
}

type Reply struct {
	Idle   time.Duration
	Mid    int64
	Repo   string
	Branch string
	Cgo    bool
}

func (r *Rpc) NewMission(args *Args, rep *Reply) error {
	log.Infof("arch: %v", args.Arch)
	log.Infof("host: %v", args.Host)
	rep.Branch = "master"
	rep.Repo = "github.com/codeskyblue/fswatch"
	rep.Mid = 1 // need to insert into mysql
	return nil
}

func (r *Rpc) UpdateBuildStatus(args *Args, reply *bool) error {
	log.Infof("update status: mid(%d) status(%s)", args.Mid, args.Status)
	*reply = true
	return nil
}

func GetMission(args *Args) (*Reply, error) {
	client, err := rpc.DialHTTP("tcp", DefaultServer)
	if err != nil {
		return nil, err
	}
	reply := &Reply{}
	err = client.Call("Rpc.NewMission", args, reply)
	return reply, err
}

func UpdateStatus(args *Args) error {
	client, err := rpc.DialHTTP("tcp", DefaultServer)
	if err != nil {
		return err
	}
	var reply bool = false
	err = client.Call("Rpc.UpdateBuildStatus", args, &reply)
	if err != nil {
		return err
	}
	if !reply {
		err = fmt.Errorf("update status meet unknwon error")
	}
	return err
}

func HandleRpc() {
	gr := new(Rpc)
	rpc.Register(gr)
	rpc.HandleHTTP()
}
