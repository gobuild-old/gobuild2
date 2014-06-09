package config

import (
	"fmt"

	"github.com/Unknwon/com"
	"github.com/codeskyblue/go-sh"

	"code.google.com/p/gcfg"
)

var Config struct {
	Server struct {
		Domain  string `gcfg:"DOMAIN"`
		RootUrl string `gcfg:"-"`
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
	c.Server.RootUrl = fmt.Sprintf("http://%s:%d", c.Server.Domain, c.Server.Port)
	return
}

// .gobuild.yml file
type PackageConfig struct {
	Filesets struct {
		Includes []string `yaml:"includes"`
		Excludes []string `yaml:"excludes"`
	} `yaml:"filesets"`
	Settings struct {
		GoFlags   string `yaml:"goflags"`
		CGOEnable bool   `yaml"cgoenable"`
	}
}

var DefaultPcfg *PackageConfig

const RCFILE = ".gobuild.yml"

func init() {
	pcfg := &PackageConfig{}
	pcfg.Filesets.Includes = []string{"README.md", "LICENSE"}
	pcfg.Filesets.Excludes = []string{".*.go"}
	pcfg.Settings.CGOEnable = true
	pcfg.Settings.GoFlags = ""
	DefaultPcfg = pcfg
}
