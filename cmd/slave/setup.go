package slave

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/log"
)

func setUp() error {
	var err error
	var binDir = filepath.Join(SELFDIR, "bin")
	var tmpDir = filepath.Join(SELFDIR, "tmp/tmp-gopath")
	os.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	if _, err := exec.LookPath("go"); err != nil {
		log.Error("require go tool installed")
		return err
	}
	sess := sh.NewSession()
	sess.SetEnv("GOBIN", binDir)
	sess.SetEnv("GOPATH", tmpDir)
	if !sh.Test("x", GOPM) {
		os.RemoveAll(filepath.Join(binDir, "gopm"))
		defer os.RemoveAll(tmpDir)
		err = sess.Command("go", "get", "-u", "-v", "github.com/gpmgo/gopm").Run()
		if err != nil {
			log.Errorf("install gopm error: %v", err)
			return err
		}
	}
	return nil
}
