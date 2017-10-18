package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func cmdNew(*cli.Context) error {
	name := conf.New.Name
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if err := os.Mkdir(name, 0774); err != nil &&
		!strings.Contains(err.Error(), "file exists") {
		return err
	}
	files := []string{
		"command-app.go",
		"build.sh",
		"main.go",
		"variables.go",
		"app.ini",
		".gitignore",
	}
	for _, v := range files {
		if err := create(name, v); err != nil {
			return err
		}
	}
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
	return ioutil.WriteFile(filePath, content, 0644)
}
