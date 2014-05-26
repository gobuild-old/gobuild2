package main

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/log"
)

var CWD, _ = os.Getwd()

func setup() {
	var gocmd string
	var err error
	var bindir = filepath.Join(CWD, "bin")
	var tmpdir = filepath.Join(CWD, "tmp-gopath")
	os.Setenv("PATH", bindir+":"+os.Getenv("PATH"))
	if gocmd, err = exec.LookPath("go"); err != nil {
		log.Fatal("require go tool installed")
	}
	sess := sh.NewSession()
	sess.SetEnv("GOBIN", bindir)
	sess.SetEnv("GOPATH", tmpdir)
	if _, err = exec.LookPath("gopm"); err != nil {
		defer os.RemoveAll(tmpdir)
		err = sess.Command("go", "get", "-u", "-v", "github.com/gpmgo/gopm").Run()
		if err != nil {
			log.Fatalf("install gopm error: %v", err)
		}
	}
	//fmt.Println(gocmd)
}
