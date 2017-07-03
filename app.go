package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dc0d/argify"
	"github.com/gobuffalo/packr"
	"github.com/hashicorp/hcl"
	"github.com/urfave/cli"
)

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

var conf struct {
	New struct {
		Name      string `envvar:"-" usage:"required"`
		Author    string `envvar:"-" value:"N/A"`
		Copyright string `envvar:"-" value:"N/A"`
	}
}

func defaultAppNameHandler() string {
	return filepath.Base(os.Args[0])
}

func defaultConfNameHandler() string {
	fp := fmt.Sprintf("%s.conf", defaultAppNameHandler())
	if _, err := os.Stat(fp); err != nil {
		fp = "app.conf"
	}
	return fp
}

func loadHCL(ptr interface{}, filePath ...string) error {
	var fp string
	if len(filePath) > 0 {
		fp = filePath[0]
	}
	if fp == "" {
		fp = defaultConfNameHandler()
	}
	cn, err := ioutil.ReadFile(fp)
	if err != nil {
		return err
	}
	err = hcl.Unmarshal(cn, ptr)
	if err != nil {
		return err
	}

	return nil
}

func app() {
	// if err := loadHCL(&conf); err != nil {
	// 	log.Println("warn:", err)
	// 	return
	// }

	app := cli.NewApp()

	{
		app.Version = "0.1.2"
		app.Author = "dc0d"
		app.Copyright = "dc0d"
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
			Name:   "new",
			Action: cmdNew,
		}
		app.Commands = append(app.Commands, c)
	}

	argify.NewArgify().Build(app, &conf)

	if err := app.Run(os.Args); err != nil {
		log.Fatalln("error:", err)
	}
}

var (
	box packr.Box
)

func init() {
	box = packr.NewBox("./cliche-lab")
}

func cmdApp(*cli.Context) error {
	defer finit(time.Second, true)
	return nil
}

func create(appName, fileName string) error {
	var content []byte
	var err error
	if content, err = box.MustBytes(fileName); err != nil {
		return err
	}
	if conf.New.Author != "N/A" && conf.New.Copyright == "N/A" {
		conf.New.Copyright = conf.New.Author
	}
	cntnt := string(content)
	cntnt = strings.Replace(cntnt, "__appname__", conf.New.Name, -1)
	cntnt = strings.Replace(cntnt, "__author__", conf.New.Author, -1)
	cntnt = strings.Replace(cntnt, "__copyright__", conf.New.Copyright, -1)
	content = []byte(cntnt)
	wd, _ := os.Getwd() // TODO:
	filePath := filepath.Join(wd, appName, fileName)
	if _, err = os.Stat(filePath); err == nil {
		return fmt.Errorf("file already exists: %v", filePath)
	}
	if err = ioutil.WriteFile(filePath, content, 0644); err != nil {
		return err
	}
	return nil
}

func cmdNew(*cli.Context) error {
	defer finit(time.Second, true)
	name := conf.New.Name
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if err := os.Mkdir(name, 0774); err != nil &&
		!strings.Contains(err.Error(), "file exists") {
		return err
	}
	files := []string{
		"app.go",
		"build.sh",
		"main.go",
		"utils.go",
		"app.conf",
		".gitignore",
	}
	for _, v := range files {
		if err := create(name, v); err != nil {
			return err
		}
	}
	return nil
}
