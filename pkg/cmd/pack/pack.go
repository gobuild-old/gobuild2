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

func findFiles(path string, depth int, focuses, skips []*regexp.Regexp) ([]string, error) {
	baseNumSeps := strings.Count(path, string(os.PathSeparator))
	var files []string
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			pathDepth := strings.Count(path, string(os.PathSeparator)) - baseNumSeps
			if pathDepth > depth {
				return filepath.SkipDir
			}
		}
		name := path //info.Name()
		isSkip := true
		for _, focus := range focuses {
			if focus.MatchString(name) {
				isSkip = false
				break
			}
		}
		for _, skip := range skips {
			if skip.MatchString(name) {
				isSkip = true
				break
			}
		}
		// log.Println(isSkip, name)
		if !isSkip {
			files = append(files, path)
			log.Debug("add file:", name, path)
		}
		return nil
	})
	return files, err
}

func init() {
	log.SetOutputLevel(log.Ldebug)
}

/*
package a program
download source
parse yaml
build binary
*/
func Action(c *cli.Context) {
	var goos, goarch = c.String("os"), c.String("arch")
	var depth = c.Int("depth")
	var output = c.String("output")
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
	log.Println("config:", pcfg)
	var focuses, skips []*regexp.Regexp
	for _, str := range pcfg.Filesets.Includes {
		focuses = append(focuses, regexp.MustCompile(str))
	}
	for _, str := range pcfg.Filesets.Excludes {
		skips = append(skips, regexp.MustCompile(str))
	}

	files, err := findFiles(".", depth, focuses, skips)
	if err != nil {
		return
	}
	log.Info("files:", files)

	// build source
	if err = sess.Command("go", "build").Run(); err != nil {
		return
	}
	cwd, _ := os.Getwd()
	program := filepath.Base(cwd)
	files = append(files, program)

	hasExt := func(ext string) bool {
		return strings.HasSuffix(output, ext)
	}

	var z Archiver
	switch {
	case hasExt(".zip"):
		fmt.Println("zip format")
		z, err = CreateZip(output)
	case hasExt(".tar"):
		fmt.Println("tar format")
		z, err = CreateTar(output)
	default:
		fmt.Println("unsupport file archive format")
		os.Exit(1)
	}
	if err != nil {
		return
	}
	log.Debug("add files")
	for _, file := range files {
		if err = z.Add(file); err != nil {
			return
		}
	}
	log.Debug("finish write zip file")
	err = z.Close()
}
