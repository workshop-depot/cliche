package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dc0d/argify"
	"github.com/dc0d/club/clubaux"
	"github.com/urfave/cli"
)

func main() {
	if err := clubaux.LoadHCL(&conf); err != nil {
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
	}

	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}
