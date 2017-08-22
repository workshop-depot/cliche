package main

import (
	"log"

	"github.com/urfave/cli"
)

var conf struct {
	Info string `envvar:"APP_INFO" usage:"sample app info" value:"bare app structure"`

	Sample struct {
		SubCommand struct {
			Param string `envvar:"-"`
		}
	}
}

func cmdApp(*cli.Context) error {
	log.Println(conf.Info, "ʕ⚆ϖ⚆ʔ")
	return nil
}
