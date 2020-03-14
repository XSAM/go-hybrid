package main

import (
	"fmt"

	"github.com/XSAM/go-hybrid/metadata"
)

func main() {
	// Set app name
	metadata.SetAppName("example")

	fmt.Printf("App info: %+v", metadata.AppInfo())
}
