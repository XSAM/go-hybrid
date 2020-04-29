package builtinutil

import (
	"errors"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestUserHomeDir(t *testing.T) {
	testCases := []struct {
		name       string
		osLibError bool
	}{
		{
			name: "os.UserHomeDir() works fine",
		},
		{
			name:       "os.UserHomeDir() return error",
			osLibError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.osLibError {
				monkey.Patch(os.UserHomeDir, func() (string, error) {
					return "", errors.New("testing")
				})
				defer func() {
					err := recover()
					assert.NotNil(t, err)
				}()
			}
			result := UserHomeDir()

			if tc.osLibError {
				// Should not run to this way
				t.Fatalf("expect panic")
			} else {
				libResult, err := os.UserHomeDir()
				assert.NoError(t, err)

				assert.Equal(t, libResult, result)
			}
		})
	}
}
