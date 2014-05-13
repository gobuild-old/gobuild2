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
	"github.com/gobuild/gobuild2/pkg/xrpc"
	"github.com/qiniu/log"
)

func sanitizedRepoName(repo string) string {
	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-4]
	}
	if strings.HasPrefix(repo, "https://") {
		repo = repo[len("https://"):]
	}
	return repo
}

var TMPDIR = "./tmp"
var PROGRAM, _ = filepath.Abs(os.Args[0])

func work(m *Mission) (err error) {
	sess := sh.NewSession()
	var gopath, _ = filepath.Abs(TMPDIR)
	sess.SetEnv("GOPATH", gopath)

	var repoAddr = m.Repo
	var cleanName = sanitizedRepoName(repoAddr)

	var srcPath = filepath.Join(gopath, "src", cleanName)
	err = sess.Command("gopm", "get", "-v", repoAddr).Run()
	if err != nil {
		log.Error(err)
		return
	}
	// TODO: change to right branch
	var outFile = "output.tar.gz"
	err = sess.Command(PROGRAM, "pack", "-o", outFile, "-gom", "gopm", sh.Dir(srcPath)).Run()
	if err != nil {
		log.Error(err)
		return
	}
	checkError := func(err error) {
		if err != nil {
			log.Errorf("err: %v", err)
		}
	}
	err = xrpc.UpdateStatus(&xrpc.Args{Mid: 1, Status: xrpc.ST_PUBLISHING})
	checkError(err)

	return nil
}

func Action(c *cli.Context) {
	fmt.Println("this is slave daemon")

	var err error
	TMPDIR, err = filepath.Abs(TMPDIR)

	if err != nil {
		log.Errorf("tmpdir to abspath err: %v", err)
		return
	}
	if !sh.Test("dir", TMPDIR) {
		os.MkdirAll(TMPDIR, 0755)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("hostname retrive err: %v", err)
	}
	args := &xrpc.Args{Os: runtime.GOOS, Arch: runtime.GOARCH, Host: hostname}
	xrpc.DefaultServer = "localhost:8010"
	for {
		reply, err := xrpc.GetMission(args)
		if err != nil {
			log.Errorf("call server rpc error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}
		log.Infof("reply: %v", reply)
		if reply.Idle != 0 {
			log.Infof("Idle for next reply: %v", reply.Idle)
			time.Sleep(reply.Idle)
		}
		missionQueue <- Mission{Repo: reply.Repo, Branch: reply.Branch, Cgo: reply.Cgo}
	}
}
