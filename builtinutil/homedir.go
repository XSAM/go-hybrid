package builtinutil

import "os"

// UserHomeDir returns the current user's home directory.
// Panic if `os.UserHomeDir()` return an error.
func UserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}
