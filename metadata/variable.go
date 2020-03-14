package metadata

import (
	"fmt"
	"runtime"

	"github.com/google/uuid"
)

var appInfo Info

func init() {
	appInfo = Info{
		RuntimeID: uuid.New().String(),
		Version: Version{
			GitVersion:   gitVersion,
			GitCommit:    gitCommit,
			GitBranch:    gitBranch,
			GitTreeState: gitTreeState,
			BuildTime:    buildTime,
			GoVersion:    runtime.Version(),
			Compiler:     runtime.Compiler,
			Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		},
	}
}
