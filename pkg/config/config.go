package config

import (
	"fmt"
	"io/ioutil"

	"code.google.com/p/gcfg"
	"github.com/Unknwon/com"
	"github.com/codeskyblue/go-sh"
	"github.com/gobuild/goyaml"
)

var Config struct {
	Server struct {
		Domain  string `gcfg:"DOMAIN"`
		RootUrl string `gcfg:"ROOTURL"`
		Addr    string `gcfg:"ADDR"`
		Port    int    `gcfg:"PORT"`
	}
	Database struct {
		DbType   string `gcfg:"DBTYPE"`
		Host     string `gcfg:"HOST"`
		Port     string `gcfg:"PORT"`
		Name     string `gcfg:"NAME"`
		User     string `gcfg:"USER"`
		Password string `gcfg:"PASSWD"`
		SslMode  string `gcfg:"SSLMODE"`
	}
	Cdn struct {
		AccessKey string `gcfg:"ACCESSKEY"`
		SecretKey string `gcfg":SECRETKEY"`
		Bulket    string `gcfg:"BULKET"`
	}
	Social map[string]*struct {
		ClientId     string `gcfg:"ID"`
		ClientSecret string `gcfg:"SECRET"`
		AuthURL      string `gcfg:"AUTHURL"`
		TokenURL     string `gcfg:"TOKENURL"`
	}
}

func Load(cfgPath string) (err error) {
	c := &Config
	if !sh.Test("file", cfgPath) {
		com.Copy(cfgPath+".default", cfgPath)
	}
	if err = gcfg.ReadFileInto(c, cfgPath); err != nil {
		return err
	}
	if c.Server.RootUrl == "" {
		c.Server.RootUrl = fmt.Sprintf("http://%s:%d", c.Server.Domain, c.Server.Port)
	}
	return
}

// .gobuild.yml file
type PackageConfig struct {
	Filesets struct {
		Includes []string `yaml:"includes"`
		Excludes []string `yaml:"excludes"`
	} `yaml:"filesets"`
	Settings struct {
		TargetDir string `yaml:"targetdir"` // target dir
		Addopts   string `yaml:"addopts"`   // extra command line options
		CGOEnable *bool  `yaml"cgoenable"`
	} `yaml:"settings"`
}

var DefaultPcfg *PackageConfig

const RCFILE = ".gobuild.yml"

func init() {
	pcfg := &PackageConfig{}
	pcfg.Filesets.Includes = []string{"README.md", "LICENSE"}
	pcfg.Filesets.Excludes = []string{".*.go"}
	// pcfg.Settings.CGOEnable = true // the default CGOEnable should be nil
	pcfg.Settings.TargetDir = ""
	pcfg.Settings.Addopts = ""
	DefaultPcfg = pcfg
}

// parse yaml
func ReadPkgConfig(filepath string) (pcfg PackageConfig, err error) {
	pcfg = PackageConfig{}
	if sh.Test("file", filepath) {
		data, er := ioutil.ReadFile(filepath)
		if er != nil {
			err = er
			return
		}
		if err = goyaml.Unmarshal(data, &pcfg); err != nil {
			return
		}
	} else {
		pcfg = *DefaultPcfg
	}
	return pcfg, nil
}
