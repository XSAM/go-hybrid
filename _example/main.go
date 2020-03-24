package main

import (
	"github.com/XSAM/go-hybrid/_example/runtime"
	"github.com/XSAM/go-hybrid/metadata"
)

func main() {
	// Set app name
	metadata.SetAppName("example")

	runtime.Start()
}
