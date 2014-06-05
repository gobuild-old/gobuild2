package slave

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/codegangsta/cli"
	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/gobuild2/models"
	"github.com/gobuild/gobuild2/pkg/base"
	"github.com/gobuild/gobuild2/pkg/xrpc"
	"github.com/gobuild/log"
)

var (
	TMPDIR     = "./tmp"
	PROGRAM, _ = filepath.Abs(os.Args[0])
	SELFDIR    = filepath.Dir(PROGRAM)
	GOPM       = filepath.Join(SELFDIR, "bin/gopm")
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

func GoInterval(dur time.Duration, f func()) chan bool {
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

func steps(m *xrpc.Mission, gopath string, sess *sh.Session, buffer *bytes.Buffer, bi xrpc.BuildInfo) (err error) {
	var task_status string
	notify := func(output string) {
		err := reportProgress(m.Mid, task_status, output)
		checkError(err)
	}
	newNotify := func(buf *bytes.Buffer) chan bool {
		return GoInterval(time.Second*2, func() {
			notify(string(buf.Bytes()))
		})
	}
	defer func() {
		fmt.Println("work steps DONE", err)
		if err != nil {
			task_status = models.ST_ERROR
			notify(err.Error())
		}
	}()

	var repoName = m.Repo
	var binName = filepath.Base(m.Repo)
	var srcPath = filepath.Join(gopath, "src", repoName)
	// var buffer = bytes.NewBuffer(nil)
	var done chan bool
	var outFile string
	var storage Storager

	if bi.UploadType == xrpc.UT_QINIU {
		var qinfo xrpc.QiniuInfo
		if err = base.Str2Objc(bi.UploadData, &qinfo); err != nil {
			return
		}
		outFile = filepath.Base(qinfo.Key)
		storage = &Qiniu{qinfo.Token, qinfo.Key, qinfo.Bulket}
	} else {
		err = fmt.Errorf("unsupported upload type:%v", bi.UploadType)
		return
	}

	var outFullPath string
	if bi.Action == models.AC_BUILD {
		outFullPath = filepath.Join(srcPath, outFile)
		task_status = models.ST_BUILDING
		done = newNotify(buffer)
		build := func() error {
			return sess.Command("go", "build", "-v", sh.Dir(srcPath)).Run()
		}
		err = build()
		done <- true
		notify(string(buffer.Bytes()))
		if err != nil {
			log.Errorf("build error: %v", err)
			return
		}
		buffer.Reset()

		// write extra pkginfo
		task_status = models.ST_PACKING
		pkginfo := "pkginfo.json"
		if err = ioutil.WriteFile(filepath.Join(srcPath, pkginfo), m.PkgInfo, 0644); err != nil {
			return
		}
		defer os.Remove(filepath.Join(srcPath, pkginfo))

		if bi.Os == "windows" {
			binName += ".exe"
		}
		err = sess.Command(PROGRAM, "pack",
			"--nobuild", "-a", pkginfo, "-a", binName, "-o", outFile, sh.Dir(srcPath)).Run()
		notify(string(buffer.Bytes()))
		if err != nil {
			log.Error(err)
			return
		}
		defer os.Remove(outFullPath)
	} else if bi.Action == models.AC_SRCPKG {
		gbcnf := `filesets:
  includes:
    - src
  excludes:
    - \.git
`
		ioutil.WriteFile(filepath.Join(gopath, ".gobuild.yml"), []byte(gbcnf), 0644)

		task_status = models.ST_PACKING
		// maybe 7 depth is enough, the hell seven
		err = sess.Command(PROGRAM, "pack", "--depth", "7", "--nobuild", "-o", "src.zip", sh.Dir(gopath)).Run()
		notify(string(buffer.Bytes()))
		if err != nil {
			return
		}
		outFullPath = filepath.Join(gopath, "src.zip")
	} else {
		err = fmt.Errorf("unknown action: %v", bi.Action)
		return
	}

	// upload and share
	task_status = models.ST_PUBLISHING
	notify(outFullPath)
	var pubAddr string
	if pubAddr, err = storage.Upload(outFullPath); err != nil {
		checkError(err)
		return
	}
	log.Debugf("publish %s to %s", outFile, pubAddr)

	reportProgress(m.Mid, models.ST_DONE, "published to "+pubAddr)
	reportPubAddr(m.Mid, pubAddr)
	return nil
}

func reportProgress(mid int64, status string, output string) error {
	log.Debugf("mid(%d) report progress, status(%s)", mid, status)
	mstatus := &xrpc.MissionStatus{
		Mid:    mid,
		Status: status,
		Output: output,
	}
	ok := false
	err := xrpc.Call("UpdateMissionStatus", mstatus, &ok)
	checkError(err)
	return err
}

func reportPubAddr(mid int64, zipballurl string) error {
	pubinfo := &xrpc.PublishInfo{
		Mid:        mid,
		ZipBallURL: zipballurl,
	}
	ok := false
	err := xrpc.Call("UpdatePubAddr", pubinfo, &ok)
	checkError(err)
	return err
}

func work(m *xrpc.Mission) (err error) {
	// notify := func(status string, output string, extra ...string) {
	// 	mstatus := &xrpc.MissionStatus{Mid: m.Mid, Status: status,
	// 		Output: output,
	// 		Extra:  strings.Join(extra, ""),
	// 	}
	// 	ok := false
	// 	err := xrpc.Call("UpdateMissionStatus", mstatus, &ok)
	// 	checkError(err)
	// }
	defer func() {
		fmt.Println("DONE", err)
		if err != nil {
			reportProgress(m.Mid, models.ST_ERROR, err.Error())
		}
	}()
	// prepare shell session
	sess := sh.NewSession()
	buffer := bytes.NewBuffer(nil)
	sess.Stdout = io.MultiWriter(buffer, os.Stdout)
	sess.Stderr = io.MultiWriter(buffer, os.Stderr)
	sess.ShowCMD = true
	gopath, err := ioutil.TempDir(TMPDIR, time.Now().Format("200601021504-"))
	if err != nil {
		log.Errorf("create gopath error: %v", err)
		return
	}
	// fmt.Println(gopath)
	// return
	// var gopath, _ = filepath.Abs(TMPDIR)
	log.Debugf("use temp gopath: %s", gopath)
	if !sh.Test("dir", gopath) {
		os.MkdirAll(gopath, 0755)
	}
	defer os.RemoveAll(gopath)
	sess.SetEnv("GOPATH", gopath)
	sess.SetEnv("CGO_ENABLE", "0")
	if m.CgoEnable {
		sess.SetEnv("CGO_ENABLE", "1")
	}
	sess.SetTimeout(time.Minute * 10) // timeout in 10minutes

	var repoName = m.Repo
	var srcPath = filepath.Join(gopath, "src", repoName)

	getsrc := func() (err error) {
		var params []interface{}
		params = append(params, "get", "-d", "-v", "-g") // todo: add -d when gopm released
		params = append(params, repoName+"@"+m.PushURI)
		// if m.Sha != "" {
		// params = append(params, repoName+"@commit:"+m.Sha)
		// } else {
		// params = append(params, repoName+"@branch:"+m.Branch)
		// }
		params = append(params, sh.Dir(gopath))
		if err = sess.Command(GOPM, params...).Run(); err != nil {
			return
		}
		if err = sess.Command("go", "get", "-v", sh.Dir(srcPath)).Run(); err != nil {
			return
		}
		return nil
	}

	newNotify := func(status string, buf *bytes.Buffer) chan bool {
		return GoInterval(time.Second*2, func() {
			reportProgress(m.Mid, status, string(buf.Bytes()))
		})
	}
	reportProgress(m.Mid, models.ST_RETRIVING, "start get source code")
	var done chan bool
	done = newNotify(models.ST_RETRIVING, buffer)
	err = getsrc()
	done <- true
	reportProgress(m.Mid, models.ST_RETRIVING, string(buffer.Bytes()))
	if err != nil {
		log.Errorf("getsource err: %v", err)
		return
	}
	buffer.Reset()
	for _, bi := range m.Builds {
		sess.SetEnv("GOOS", bi.Os)
		sess.SetEnv("GOARCH", bi.Arch)
		steps(m, gopath, sess, buffer, bi)
	}
	return nil

	// var outFile = filepath.Base(m.UpKey)
	// var outFullPath = filepath.Join(srcPath, outFile)

	// done = newNotify(models.ST_BUILDING, buffer)

	// err = sess.Command(GOPM, "build", "-u", "-v", sh.Dir(srcPath)).Run()
	// err = build()
	// done <- true
	// notify(models.ST_BUILDING, string(buffer.Bytes()))
	// if err != nil {
	// 	log.Errorf("build error: %v", err)
	// 	return
	// }
	// buffer.Reset()

	// write extra pkginfo
	// pkginfo := "pkginfo.json"
	// ioutil.WriteFile(filepath.Join(srcPath, pkginfo), m.PkgInfo, 0644)

	// err = sess.Command(PROGRAM, "pack",
	// 	"--nobuild", "-a", pkginfo, "-o", outFile, sh.Dir(srcPath)).Run()
	// notify(models.ST_PACKING, string(buffer.Bytes()))
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }

	// var cdnPath = m.UpKey
	// notify(models.ST_PUBLISHING, cdnPath)
	// log.Infof("cdn path: %s", cdnPath)
	// q := &Qiniu{m.UpToken, m.UpKey, m.Bulket} // uptoken, key}
	// var pubAddr string
	// if pubAddr, err = q.Upload(outFullPath); err != nil {
	// 	checkError(err)
	// 	return
	// }

	// log.Debugf("publish %s to %s", outFile, pubAddr)
	// notify(models.ST_DONE, pubAddr)
	// return nil
}

func init() {
	var err error
	HOSTNAME, err = os.Hostname()
	if err != nil {
		log.Fatalf("hostname retrive err: %v", err)
	}
}

var IsPrivateUpload bool //todo

func prepare() (err error) {
	TMPDIR, err = filepath.Abs(TMPDIR)
	if err != nil {
		log.Errorf("tmpdir to abspath err: %v", err)
		return
	}
	if !sh.Test("dir", TMPDIR) {
		os.MkdirAll(TMPDIR, 0755)
	}
	if err = setUp(); err != nil {
		log.Fatalf("setUp environment error:%v", err)
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

		if mission.Idle != 0 {
			fmt.Print(".")
			time.Sleep(mission.Idle)
			continue
		}
		log.Infof("new mission from xrpc: %s", base.Objc2Str(mission))
		missionQueue <- mission
	}
}
