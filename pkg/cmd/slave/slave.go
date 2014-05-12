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
	"github.com/gobuild/gobuild2/routers"
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

func work(m *Mission) {
	sess := sh.NewSession()
	sess.SetEnv("GOPATH", TMPDIR)

	var err error
	var repoAddr = m.Repo
	var cleanName = sanitizedRepoName(repoAddr)

	//var repoAddr = "github.com/codeskyblue/fswatch"
	var srcPath = filepath.Join(TMPDIR, "src", cleanName)
	_ = srcPath
	// os.MkdirAll(srcPath, 0755)
	err = sess.Command("gopm", "get", "-v", repoAddr).Run() //, sh.Dir(filepath.Dir(srcPath))).Run()
	if err != nil {
		log.Error(err)
		return
	}
	// TODO: change to right branch
	var PROGRAM, _ = filepath.Abs(os.Args[0])
	fmt.Println(PROGRAM)
	err = sess.Command(PROGRAM, "pack", "-o", "output.tar.gz", "-gom", "gopm", sh.Dir(srcPath)).Run()
	if err != nil {
		log.Error(err)
	}
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
	args := &routers.Args{Os: runtime.GOOS, Arch: runtime.GOARCH, Host: hostname}
	for {
		reply, err := routers.GetMission("localhost:8010", args)
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
