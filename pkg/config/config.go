package config

import (
	"fmt"

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
}

func Load(cfgPath string) (err error) {
	c := &Config
	err = gcfg.ReadFileInto(c, cfgPath)
	c.Server.RootUrl = fmt.Sprintf("http://%s:%d", c.Server.Domain, c.Server.Port)
	return
}
