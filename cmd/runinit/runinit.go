package runinit

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/codegangsta/cli"
	"github.com/gobuild/gobuild2/pkg/config"
	"github.com/gobuild/goyaml"
)

func Action(c *cli.Context) {
	// content, _ := ioutil.ReadFile(config.RCFILE)
	// goyaml.Unmarshal(content, &config.DefaultPcfg)
	// fmt.Println(string(content))
	// fmt.Println(config.DefaultPcfg)

	data, _ := goyaml.Marshal(config.DefaultPcfg)
	if err := ioutil.WriteFile(config.RCFILE, data, 0644); err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("rcfile: %s created", config.RCFILE)
	}
}
