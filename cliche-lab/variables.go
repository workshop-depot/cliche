package main

//-----------------------------------------------------------------------------

var conf struct {
	Info string `envvar:"APP_INFO" usage:"sample app info" value:"bare app structure"`

	Sample struct {
		SubCommand struct {
			Param string `envvar:"-"`
		}
	}
}

//-----------------------------------------------------------------------------

// build flags
var (
	BuildTime  string
	CommitHash string
	GoVersion  string
	GitTag     string
)

//-----------------------------------------------------------------------------
