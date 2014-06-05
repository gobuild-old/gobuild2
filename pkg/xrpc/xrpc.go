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
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/gobuild/log"
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

var DefaultWebAddress = "localhost:8010"

type Rpc struct{}

type HostInfo struct {
	Os, Arch string
	Host     string
}

type QiniuInfo struct {
	Key    string
	Token  string
	Bulket string
}

type PublishInfo struct {
	Mid        int64
	ZipBallUrl string
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

	Repo    string
	PushURI string

	CgoEnable bool

	PkgInfo []byte // store it to gobuild.pkginfo
	Builds  []BuildInfo
}

const (
	UT_QINIU = "up-qiniu"
	UT_FTP   = "up-ftp"
)

type BuildInfo struct {
	Action     string // build or just package source
	CgoEnable  bool
	Os, Arch   string
	UploadType string // cnd(like qiniu) or ftp
	UploadData string // encoded data
}

func Call(method string, args interface{}, reply interface{}) error {
	client, err := rpc.DialHTTP("tcp", DefaultWebAddress)
	if err != nil {
		return err
	}
	defer client.Close()
	return client.Call("Rpc."+method, args, reply)
}

var defaultBulket string

// generate qiniu token
func qntoken(key string) string {
	scope := defaultBulket + ":" + key
	log.Infof("qiniu scrope: %s", scope)
	policy := rs.PutPolicy{
		Expires: uint32(time.Now().Unix() + 3600),
		Scope:   scope,
	}
	return policy.Token(nil)
}

func (r *Rpc) GetMission(args *HostInfo, rep *Mission) error {
	log.Debugf("arch: %v, host: %v", args.Arch, args.Host)
	tasks, err := models.GetAvaliableTasks(args.Os, args.Arch)
	if err == models.ErrTaskNotAvaliable {
		rep.Idle = time.Second * 3
		return nil
	}
	if err != nil {
		log.Errorf("rpc: get mission error: %v", err)
		return err
	}

	task := tasks[0] // use first task
	rep.Mid = task.Id
	rep.Repo = task.Repo.Uri
	rep.PushURI = task.PushType + ":" + task.PushValue
	rep.CgoEnable = task.CgoEnable
	rep.PkgInfo, _ = json.MarshalIndent(PkgInfo{
		PushURI:     task.PushType + ":" + task.PushValue,
		Author:      []string{"unknown"},
		Description: "unknown",
	}, "", "    ")

	for _, tk := range tasks {
		if tk.TagBranch == "" {
			tk.TagBranch = "temp-" + tk.PushType + ":" + tk.PushValue
		}
		filename := fmt.Sprintf("%s-%s-%s.%s", filepath.Base(rep.Repo), tk.Os, tk.Arch, "zip")
		if tk.Action == models.AC_SRCPKG {
			filename = fmt.Sprintf("%s-all-source.%s", filepath.Base(rep.Repo), "zip")
		}
		key := com.Expand("m{tid}/{reponame}/br-{branch}/{filename}", map[string]string{
			"tid":      strconv.Itoa(int(rep.Mid)),
			"reponame": rep.Repo,
			"branch":   tk.TagBranch,
			"filename": filename,
		})
		bi := BuildInfo{
			Action:     tk.Action,
			Os:         tk.Os,
			Arch:       tk.Arch,
			UploadType: UT_QINIU,
			UploadData: base.Objc2Str(QiniuInfo{
				Bulket: defaultBulket,
				Key:    key,
				Token:  qntoken(key),
			}),
		}
		rep.Builds = append(rep.Builds, bi)
	}
	return nil
}

type PkgInfo struct {
	MainFile    string   `json:"main_file"`
	Author      []string `json:"author"`
	From        string   `json:"from"`
	Description string   `json:"description"`
	Created     string   `json:"created"`
	Os          string   `json:"os"`
	Arch        string   `json:"arch"`
	PushURI     string   `json:"push_uri"`
}

func (r *Rpc) UpdatePubAddr(args *PublishInfo, reply *bool) error {
	log.Infof("pub addr %v", *args)
	*reply = true
	err := models.UpdatePubAddr(args.Mid, args.ZipBallUrl)
	return err

}

func (r *Rpc) UpdateMissionStatus(args *MissionStatus, reply *bool) error {
	log.Debugf("update status: mid(%d) status(%s) extra(%s)", args.Mid, args.Status, args.Extra)
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
