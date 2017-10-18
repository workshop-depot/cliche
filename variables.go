package main

import "github.com/gobuffalo/packr"

//-----------------------------------------------------------------------------

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

//-----------------------------------------------------------------------------

var conf struct {
	New struct {
		Name      string `envvar:"-" usage:"required"`
		Author    string `envvar:"-" value:"N/A"`
		Copyright string `envvar:"-" value:"N/A"`
	}
}

//-----------------------------------------------------------------------------

var (
	box packr.Box
)

func init() {
	box = packr.NewBox("./cliche-lab")
}

//-----------------------------------------------------------------------------
