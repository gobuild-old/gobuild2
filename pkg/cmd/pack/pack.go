package pack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/gobuild/log"
	"launchpad.net/goyaml"

	"github.com/codegangsta/cli"
	sh "github.com/codeskyblue/go-sh"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/unknwon/com"
)

func init() {
	log.SetFlags(log.Linfo)
	// log.SetOutputLevel(log.Ldebug)
}

func findFiles(path string, depth int, skips []*regexp.Regexp) ([]string, error) {
	baseNumSeps := strings.Count(path, string(os.PathSeparator))
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pathDepth := strings.Count(path, string(os.PathSeparator)) - baseNumSeps
			if pathDepth > depth {
				return filepath.SkipDir
			}
		}
		name := info.Name()
		isSkip := false
		for _, skip := range skips {
			if skip.MatchString(name) {
				isSkip = true
				break
			}
		}
		// log.Println(isSkip, name)
		if !isSkip {
			files = append(files, path)
			log.Info("add file:", name, path)
		}
		return nil
	})
	return files, err
}

func Action(c *cli.Context) {
	var goos, goarch = c.String("os"), c.String("arch")
	var depth = c.Int("depth")
	var output = c.String("output")
	var gom = c.String("gom")

	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()
	sess := sh.NewSession()
	sess.SetEnv("GOOS", goos)
	sess.SetEnv("GOARCH", goarch)
	sess.ShowCMD = true
	// parse yaml
	var pcfg = new(config.PackageConfig)
	if com.IsExist(config.RCFILE) {
		data, er := ioutil.ReadFile(config.RCFILE)
		if er != nil {
			err = er
			return
		}
		if err = goyaml.Unmarshal(data, pcfg); err != nil {
			return
		}
	} else {
		pcfg = config.DefaultPcfg
	}
	log.Debug("config:", pcfg)

	var skips []*regexp.Regexp
	for _, str := range pcfg.Filesets.Excludes {
		skips = append(skips, regexp.MustCompile("^"+str+"$"))
	}

	var files []string
	for _, filename := range pcfg.Filesets.Includes {
		fs, err := findFiles(filename, depth, skips)
		if err != nil {
			return
		}
		files = append(files, fs...)
	}

	log.Infof("archive file to: %s", output)
	var z Archiver
	hasExt := func(ext string) bool { return strings.HasSuffix(output, ext) }
	switch {
	case hasExt(".zip"):
		fmt.Println("zip format")
		z, err = CreateZip(output)
	case hasExt(".tar"):
		fmt.Println("tar format")
		z, err = CreateTar(output)
	case hasExt(".tgz"):
		fallthrough
	case hasExt(".tar.gz"):
		fmt.Println("tar.gz format")
		z, err = CreateTgz(output)
	default:
		fmt.Println("unsupport file archive format")
		os.Exit(1)
	}
	if err != nil {
		return
	}

	// build source
	if err = sess.Command(gom, "build").Run(); err != nil {
		return
	}
	cwd, _ := os.Getwd()
	program := filepath.Base(cwd)
	files = append(files, program)

	log.Debug("archive files")
	for _, file := range files {
		if err = z.Add(file); err != nil {
			return
		}
	}
	log.Info("finish archive file")
	err = z.Close()
}
