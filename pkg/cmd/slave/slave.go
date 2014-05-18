package slave

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/xrpc"
	"github.com/qiniu/log"
)

func sanitizedRepoPath(repo string) string {
	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-4]
	}
	if strings.HasPrefix(repo, "https://") {
		repo = repo[len("https://"):]
	}
	return repo
}

var (
	TMPDIR     = "./tmp"
	PROGRAM, _ = filepath.Abs(os.Args[0])
	HOSTNAME   = "localhost"
	HOSTINFO   = &xrpc.HostInfo{Os: runtime.GOOS, Arch: runtime.GOARCH, Host: HOSTNAME}
)

func checkError(err error) {
	if err != nil {
		log.Errorf("err: %v", err)
	}
}

func work(m *xrpc.Mission) (err error) {
	notify := func(status string, extra ...string) {
		mstatus := &xrpc.MissionStatus{Mid: m.Mid, Status: status, Extra: strings.Join(extra, "")}
		ok := false
		err := xrpc.Call("UpdateMissionStatus", mstatus, &ok)
		checkError(err)
	}
	defer func() {
		if err != nil {
			notify(models.ST_ERROR, err.Error())
		}
	}()
	// prepare shell session
	sess := sh.NewSession()
	sess.ShowCMD = true
	var gopath, _ = filepath.Abs(TMPDIR)
	sess.SetEnv("GOPATH", gopath)
	sess.SetEnv("CGO_ENABLE", "")
	if m.CgoEnable {
		sess.SetEnv("CGO_ENABLE", "1")
	}
	sess.SetEnv("GOOS", m.Os)
	sess.SetEnv("GOARCH", m.Arch)

	var repoAddr = m.Repo
	var cleanRepoName = sanitizedRepoPath(repoAddr)

	notify(models.ST_RETRIVING)
	var srcPath = filepath.Join(gopath, "src", cleanRepoName)

	getsrc := func() (err error) {
		if err = sess.Command("gopm", "get", "-v", "-u", repoAddr).Run(); err != nil {
			return
		}
		if err = sess.Command("gopm", "get", "-v", repoAddr).Run(); err != nil {
			return
		}
		return nil
	}

	// get source code
	if err = getsrc(); err != nil {
		log.Errorf("getsource err: %v", err)
		return
	}

	// TODO: change to right branch
	var outFile = fmt.Sprintf("%s-%s.%s", filepath.Base(cleanRepoName), m.Branch, "tar.gz")
	var outFullPath = filepath.Join(srcPath, outFile)
	notify(models.ST_BUILDING)
	err = sess.Command(PROGRAM, "pack", "-o", outFile, "-gom", "gopm", sh.Dir(srcPath)).Run()
	if err != nil {
		log.Error(err)
		return
	}
	notify(models.ST_PUBLISHING)
	// timestamp := time.Now().Format("20060102-150405")
	var cdnPath = fmt.Sprintf("m%d/%s/%s", m.Mid, cleanRepoName, outFile)
	log.Infof("cdn path: %s", cdnPath)
	var pubAddress string
	if pubAddress, err = UploadQiniu(outFullPath, cdnPath); err != nil {
		checkError(err)
		return
	}
	log.Debugf("publish %s to %s", outFile, pubAddress)
	notify(models.ST_DONE, pubAddress)
	return nil
}

func init() {
	var err error
	HOSTNAME, err = os.Hostname()
	if err != nil {
		log.Fatalf("hostname retrive err: %v", err)
	}
}

func prepare() (err error) {
	qi := new(xrpc.QiniuInfo)
	xrpc.Call("GetQiniuInfo", HOSTINFO, qi)

	initQiniu(qi.AccessKey, qi.SecretKey, qi.Bulket)

	TMPDIR, err = filepath.Abs(TMPDIR)
	if err != nil {
		log.Errorf("tmpdir to abspath err: %v", err)
		return
	}
	if !sh.Test("dir", TMPDIR) {
		os.MkdirAll(TMPDIR, 0755)
	}
	startWork()
	return nil
}

func Action(c *cli.Context) {
	fmt.Println("this is slave daemon")
	webaddr := c.String("webaddr")
	xrpc.DefaultWebAddress = webaddr

	if err := prepare(); err != nil {
		log.Fatalf("slave prepare err: %v", err)
	}
	for {
		mission := &xrpc.Mission{}
		if err := xrpc.Call("GetMission", HOSTINFO, mission); err != nil {
			log.Errorf("get mission failed: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		//log.Infof("reply: %v", mission)
		if mission.Idle != 0 {
			//log.Infof("Idle for next reply: %v", mission.Idle)
			fmt.Print(".")
			time.Sleep(mission.Idle)
			continue
		}
		log.Infof("new mission from xrpc: %v", mission)
		missionQueue <- mission
	}
}
