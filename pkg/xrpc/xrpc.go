package xrpc

import (
	"net/rpc"
	"time"

	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/qiniu/log"
)

var DefaultWebAddress = "localhost:8010"

type Rpc struct{}

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
	Idle      time.Duration
	Mid       int64
	Repo      string
	Branch    string
	CgoEnable bool
	Os, Arch  string
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

var missionQueue = make(chan *Mission, 1)

func init() {
	go func() {
		for {
			missionQueue <- &Mission{Repo: "github.com/wangwenbin/2048-go", Branch: "master", Mid: 2,
				CgoEnable: true, Os: "windows", Arch: "386"}
			time.Sleep(5 * time.Second)
		}
	}()
}

func (r *Rpc) GetMission(args *HostInfo, rep *Mission) error {
	log.Infof("arch: %v", args.Arch)
	log.Infof("host: %v", args.Host)
	task, err := models.GetAvaliableTask(args.Os, args.Arch)
	switch err {
	case nil:
		rep.CgoEnable = task.CgoEnable
		rep.Os, rep.Arch = task.Os, task.Arch
		rep.Mid = task.Id
		rep.Repo = task.Repo.Uri
		rep.Branch = task.Branch
		return nil
	case models.ErrTaskNotAvaliable:
		rep.Idle = time.Second
		return nil
	default:
		return err
	}
}

func (r *Rpc) UpdateMissionStatus(args *MissionStatus, reply *bool) error {
	log.Infof("update status: mid(%d) status(%s) extra(%s)", args.Mid, args.Status, args.Extra)
	*reply = true
	err := models.UpdateTaskStatus(args.Mid, args.Status, args.Extra)
	return err
}

func HandleRpc() {
	gr := new(Rpc)
	rpc.Register(gr)
	rpc.HandleHTTP()
}
