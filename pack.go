package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"launchpad.net/goyaml"

	"github.com/codegangsta/cli"
	sh "github.com/codeskyblue/go-sh"
	"github.com/unknwon/com"
)

func FindFiles(path string, depth int, focuses, skips []*regexp.Regexp) ([]string, error) {
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
		}
		return nil
	})
	return files, err
}

/*
package a program
download source
parse yaml
build binary
*/
func runPack(c *cli.Context) {
	var goos, goarch = c.String("os"), c.String("arch")
	var depth = c.Int("depth")
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
	var pcfg = new(PackageConfig)
	if com.IsExist(RCFILE) {
		data, er := ioutil.ReadFile(RCFILE)
		if er != nil {
			err = er
			return
		}
		if err = goyaml.Unmarshal(data, pcfg); err != nil {
			return
		}
	} else {
		pcfg = defaultDcfg
	}
	log.Println("config:", pcfg)
	var focuses, skips []*regexp.Regexp
	for _, str := range pcfg.Filesets.Includes {
		focuses = append(focuses, regexp.MustCompile(str))
	}
	for _, str := range pcfg.Filesets.Excludes {
		skips = append(skips, regexp.MustCompile(str))
	}

	files, err := FindFiles(".", depth, focuses, skips)
	if err != nil {
		return
	}
	log.Println("files:", files)

	// build source
	if err = sess.Command("go", "build").Run(); err != nil {
		return
	}
	// tar files
	// todo
}
