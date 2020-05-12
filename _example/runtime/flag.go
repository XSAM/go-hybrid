package runtime

import "time"

type Flag struct {
	Number      int           `flag:"env"`
	Duration    time.Duration `flag:""`
	Environment Environment   `flag:""`
}

type Environment struct {
	DevelopmentMode bool `flag:"env short=d" flag-usage:"change environment mode to development"`
	JSONLogStyle    bool `flag:"env name=json-log-style short=j" flag-usage:"change log style to JSON"`
}
