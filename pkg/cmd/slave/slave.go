package slave

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

func Action(c *cli.Context) {
	fmt.Println("this is slave daemon")

	var TMPDIR = "./tmp"
	var err error
	TMPDIR, err = filepath.Abs(TMPDIR)

	args := &routers.Args{Os: runtime.GOOS, Arch: runtime.GOARCH}
	reply, err := routers.GetMission("localhost:8010", args)
	log.Infof("reply: %v", reply)

	if err != nil {
		log.Errorf("tmpdir to abspath err: %v", err)
		return
	}
	if !sh.Test("dir", TMPDIR) {
		os.MkdirAll(TMPDIR, 0755)
	}
	sess := sh.NewSession()
	sess.SetEnv("GOPATH", TMPDIR)

	var repoAddr = "github.com/shxsun/fswatch"
	var cleanName = sanitizedRepoName(repoAddr)
	var srcPath = filepath.Join(TMPDIR, "src", cleanName)
	_ = srcPath
	// os.MkdirAll(srcPath, 0755)
	err = sess.Command("gopm", "get", "-v", repoAddr).Run() //, sh.Dir(filepath.Dir(srcPath))).Run()
	if err != nil {
		log.Error(err)
		return
	}
	var PROGRAM, _ = filepath.Abs(os.Args[0])
	fmt.Println(PROGRAM)
	err = sess.Command(PROGRAM, "pack", "-o", "output.tar.gz", "-gom", "gopm", sh.Dir(srcPath)).Run()
	if err != nil {
		log.Error(err)
	}
}
