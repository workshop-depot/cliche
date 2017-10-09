package main

import (
	"log"

	"github.com/urfave/cli"
)

func cmdApp(*cli.Context) error {
	log.Println(conf.Info, "ʕ⚆ϖ⚆ʔ")
	return nil
}
