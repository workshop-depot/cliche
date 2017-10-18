package main

//-----------------------------------------------------------------------------

var conf struct {
	PostgreSQL struct {
		User     string
		Password string
		Database string
	}

	NATS struct {
		NATSURL      string
		NATSUser     string
		NATSPassword string
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
