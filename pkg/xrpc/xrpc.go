package xrpc

import (
	"encoding/json"
	"fmt"
	"net/rpc"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Unknwon/com"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
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
	Output string
	Extra  string
}

type Mission struct {
	Idle time.Duration
	Mid  int64

	Repo   string
	Branch string
	Sha    string

	UpToken string // for qiniu upload
	UpKey   string // for qiniu upload
	Bulket  string

	CgoEnable bool
	Os, Arch  string

	PkgInfo []byte // store it to gobuild.pkginfo
}

func Call(method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", DefaultWebAddress)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call("Rpc."+method, args, reply)
}

type PkgInfo struct {
	MainFile    string   `json:"main_file"`
	Author      []string `json:"author"`
	From        string   `json:"from"`
	Description string   `json:"description"`
	Created     string   `json:"created"`
	Os          string   `json:"os"`
	Arch        string   `json:"arch"`
	Sha         string   `json:"sha"`
}

var defaultBulket string

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
		rep.Sha = task.Sha

		// rep.UpKey
		filename := fmt.Sprintf("%s-%s-%s.%s", filepath.Base(rep.Repo), rep.Os, rep.Arch, "zip")
		rep.UpKey = com.Expand("m{tid}/{reponame}/br-{branch}/{filename}", map[string]string{
			"tid":      strconv.Itoa(int(rep.Mid)),
			"reponame": rep.Repo,
			"branch":   rep.Branch,
			"filename": filename,
		})
		policy := rs.PutPolicy{
			Scope: defaultBulket + ":" + rep.UpKey,
		}
		policy.Expires = uint32(time.Now().Unix() + 3600)
		rep.UpToken = policy.Token(nil)
		rep.Bulket = defaultBulket

		// todo
		rep.PkgInfo, _ = json.MarshalIndent(PkgInfo{
			Sha:         task.Sha,
			Author:      []string{"unknown"},
			Description: "unknown",
		}, "", "    ")
		return nil
	case models.ErrTaskNotAvaliable:
		rep.Idle = time.Second * 3
		return nil
	default:
		return err
	}
}

func (r *Rpc) UpdateMissionStatus(args *MissionStatus, reply *bool) error {
	log.Infof("update status: mid(%d) status(%s) extra(%s)", args.Mid, args.Status, args.Extra)
	*reply = true
	err := models.UpdateTaskStatus(args.Mid, args.Status, args.Output)
	return err
}

func HandleRpc() {
	conf.ACCESS_KEY = config.Config.Cdn.AccessKey
	conf.SECRET_KEY = config.Config.Cdn.SecretKey
	defaultBulket = config.Config.Cdn.Bulket

	gr := new(Rpc)
	rpc.Register(gr)
	rpc.HandleHTTP()
}
