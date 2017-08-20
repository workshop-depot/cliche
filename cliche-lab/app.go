package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dc0d/argify"
	"github.com/dc0d/cliche/tad"
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

func app() {
	if err := tad.LoadHCL(&conf); err != nil {
		log.Println("warn:", err)
	}

	app := cli.NewApp()

	{
		app.Version = "0.0.1"
		app.Author = "__author__"
		app.Copyright = "__copyright__"
		now := time.Now()
		app.Description = fmt.Sprintf(
			"Build Time:  %v %v\n   Go:          %v\n   Commit Hash: %v\n   Git Tag:     %v",
			now.Weekday(),
			BuildTime,
			GoVersion,
			CommitHash,
			GitTag)
		app.Name = "__appname__"
		app.Usage = ""
	}

	{
		app.Action = cmdApp

		/*
			c := cli.Command{
				Name: `sample`,
			}
			c.Subcommands = append(c.Subcommands, cli.Command{
				Name:   "subcommand",
				Action: cmdSampleSubCommand,
			})
			app.Commands = append(app.Commands, c)
		*/
	}

	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}

func cmdApp(*cli.Context) error {
	defer tad.Finit(time.Second, true)
	log.Println(conf.Info, "ʕ⚆ϖ⚆ʔ")
	return nil
}

/*
func cmdSampleSubCommand(*cli.Context) error {
	defer finit(time.Second, true)
	log.Println(conf.Sample.SubCommand.Param)
	return nil
}
*/
