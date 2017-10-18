package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dc0d/argify"
	"github.com/dc0d/config/iniconfig"
	"github.com/urfave/cli"
)

func main() {
	if err := iniconfig.New().Load(&conf); err != nil {
		// this error does not help much, unless we explicitly need it
		// in which case it should be handled properly
	}
	// we can use hcl format too. the only difference is
	// the file extension which is .conf (instead of .ini)
	// if err := hclconfig.New().Load(&conf); err != nil {
	// }

	app := cli.NewApp()
	setAppInfo(app)
	addCommands(app)
	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func addCommands(app *cli.App) {
	app.Action = cmdApp

	// other commands will be added here
	// cmd1 := cli.Command{...}
	// cmd2 := cli.Command{...}
	// app.Commands = append(app.Commands, cmd1, cmd2)
}

func setAppInfo(app *cli.App) {
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
