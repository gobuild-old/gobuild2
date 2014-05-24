package slave

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/xrpc"
	"github.com/qiniu/log"
)

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

type NTMsg struct {
	Status string
	Output string
	Extra  string
}

func work(m *xrpc.Mission) (err error) {
	notify := func(status string, output string, extra ...string) {
		mstatus := &xrpc.MissionStatus{Mid: m.Mid, Status: status,
			Output: output,
			Extra:  strings.Join(extra, ""),
		}
		ok := false
		err := xrpc.Call("UpdateMissionStatus", mstatus, &ok)
		checkError(err)
	}
	defer func() {
		fmt.Println("DONE", err)
		if err != nil {
			notify(models.ST_ERROR, err.Error())
		}
	}()
	// prepare shell session
	sess := sh.NewSession()
	buffer := bytes.NewBuffer(nil)
	sess.Stdout = io.MultiWriter(buffer, os.Stdout)
	sess.Stderr = io.MultiWriter(buffer, os.Stderr)
	sess.ShowCMD = true
	var gopath, _ = filepath.Abs(TMPDIR)
	sess.SetEnv("GOPATH", gopath)
	sess.SetEnv("CGO_ENABLE", "0")
	if m.CgoEnable {
		sess.SetEnv("CGO_ENABLE", "1")
	}
	sess.SetEnv("GOOS", m.Os)
	sess.SetEnv("GOARCH", m.Arch)

	var repoAddr = m.Repo
	var cleanRepoName = base.SanitizedRepoPath(repoAddr)

	var srcPath = filepath.Join(gopath, "src", cleanRepoName)

	getsrc := func() (err error) {
		// if err = sess.Command("gopm", "get", "-v", "-u", repoAddr).Run(); err != nil {
		// return
		// }
		if err = sess.Command("gopm", "get", "-g", "-v", repoAddr).Run(); err != nil {
			// if err = sess.Command("gopm", "get", "-v", repoAddr+"@commit:"+m.Sha).Run(); err != nil {
			return
		}
		return nil
	}

	GoInterval := func(dur time.Duration, f func()) chan bool {
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case <-time.After(dur):
					f()
				}
			}
		}()
		return done
	}

	newNotify := func(status string, buf *bytes.Buffer) chan bool {
		return GoInterval(time.Second*2, func() {
			notify(status, string(buf.Bytes()))
		})
	}

	notify(models.ST_RETRIVING, "start get source")
	var done chan bool
	done = newNotify(models.ST_RETRIVING, buffer)
	err = getsrc()
	done <- true
	notify(models.ST_RETRIVING, string(buffer.Bytes()))
	if err != nil {
		log.Errorf("getsource err: %v", err)
		return
	}
	buffer.Reset()

	extention := "zip"
	var outFile = fmt.Sprintf("%s-%s-%s.%s", filepath.Base(cleanRepoName), m.Os, m.Arch, extention)
	var outFullPath = filepath.Join(srcPath, outFile)

	// notify(models.ST_BUILDING, "start building")
	done = newNotify(models.ST_BUILDING, buffer)
	err = sess.Command("gopm", "build", "-u", "-v", sh.Dir(srcPath)).Run()
	done <- true
	notify(models.ST_BUILDING, string(buffer.Bytes()))
	if err != nil {
		log.Errorf("gopm build error: %v", err)
		return
	}
	buffer.Reset()

	err = sess.Command(PROGRAM, "pack", "--nobuild", "-o", outFile, sh.Dir(srcPath)).Run()
	notify(models.ST_PACKING, string(buffer.Bytes()))
	if err != nil {
		log.Error(err)
		return
	}

	notify(models.ST_PUBLISHING, "")
	// timestamp := time.Now().Format("20060102-150405")
	var cdnPath = fmt.Sprintf("m%d/%s/raw/%s", m.Mid, cleanRepoName, outFile)
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
