package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

//-----------------------------------------------------------------------------

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

// LoadHCL loads hcl conf file. default conf file names (if filePath not provided)
// in the same directory are <appname>.conf and if not fount app.conf
func LoadHCL(ptr interface{}, filePath ...string) error {
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

//-----------------------------------------------------------------------------
