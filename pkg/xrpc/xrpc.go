package xrpc

import (
	"net/rpc"
	"time"

	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/qiniu/log"
)

var DefaultWebAddress = "localhost:8010"

type Rpc struct{}

const (
	ST_PENDING    = "pending"
	ST_RETRIVING  = "retriving"
	ST_BUILDING   = "building"
	ST_PUBLISHING = "publishing"
	ST_DONE       = "done"
	ST_ERROR      = "error"
)

type HostInfo struct {
	Os, Arch string
	Host     string
}

type QiniuInfo struct {
	AccessKey string
	SecretKey string
	Bulket    string
}

type MissionStatus struct {
	Mid    int64 // mission id
	Status string
	Extra  string
}

type Mission struct {
	Idle   time.Duration
	Mid    int64
	Repo   string
	Branch string
	Cgo    bool
}

func (r *Rpc) GetQiniuInfo(args *HostInfo, rep *QiniuInfo) error {
	log.Infof("arch: %v", args.Arch)
	log.Infof("host: %v", args.Host)
	cdn := config.Config.Cdn
	rep.AccessKey = cdn.AccessKey
	rep.SecretKey = cdn.SecretKey
	rep.Bulket = cdn.Bulket
	return nil
}

func Call(method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", DefaultWebAddress)
	if err != nil {
		return err
	}
	return client.Call("Rpc."+method, args, reply)
}

func (r *Rpc) GetMission(args *HostInfo, rep *Mission) error {
	log.Infof("arch: %v", args.Arch)
	log.Infof("host: %v", args.Host)
	rep.Branch = "master"
	rep.Repo = "github.com/codeskyblue/fswatch"
	rep.Mid = 1 // need to insert into mysql
	return nil
}

func (r *Rpc) UpdateMissionStatus(args *MissionStatus, reply *bool) error {
	log.Infof("update status: mid(%d) status(%s) extra(%s)", args.Mid, args.Status, args.Extra)
	*reply = true
	return nil
}

func HandleRpc() {
	gr := new(Rpc)
	rpc.Register(gr)
	rpc.HandleHTTP()
}
