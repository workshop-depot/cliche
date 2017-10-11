package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dc0d/argify"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	{
		app.Version = "0.11.2"
		app.Author = "dc0d"
		app.Copyright = "kaveh.shahbazian@gmail.com"
		now := time.Now()
		app.Description = fmt.Sprintf(
			"Build Time:  %v %v\n   Go:          %v\n   Commit Hash: %v\n   Git Tag:     %v",
			now.Weekday(),
			BuildTime,
			GoVersion,
			CommitHash,
			GitTag)
		app.Name = "cliche"
		app.Usage = ""
	}

	{
		c := cli.Command{
			Name:        "new",
			Action:      cmdNew,
			Usage:       "cliche new --name <cli_app_name> --author <author> --copyright <copyright>",
			Description: "creates a new cli app",
		}
		app.Commands = append(app.Commands, c)
	}

	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}
