package main

import (
	"fmt"

	"code.google.com/p/gcfg"
)

type Config struct {
	Server struct {
		Domain  string `gcfg:"DOMAIN"`
		RootUrl string `gcfg:"-"`
		Addr    string `gcfg:"ADDR"`
		Port    int    `gcfg:"PORT"`
	}
	Database struct {
		DbType   string `gcfg:"DBTYPE"`
		Host     string `gcfg:"HOST"`
		Name     string `gcfg:"NAME"`
		User     string `gcfg:"USER"`
		Password string `gcfg:"PASSWD"`
	}
}

func readCfg(cfgPath string) (c *Config, err error) {
	c = new(Config)
	err = gcfg.ReadFileInto(c, cfgPath)
	c.Server.RootUrl = fmt.Sprintf("http://%s:%d", c.Server.Domain, c.Server.Port)
	return
}
