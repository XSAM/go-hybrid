package cmdutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_resolveFlagTag(t *testing.T) {
	testCases := []struct {
		structTag       string
		expectedFlagTag flagTag
	}{
		{
			structTag:       "",
			expectedFlagTag: flagTag{enable: false},
		},
		{
			structTag:       "foo",
			expectedFlagTag: flagTag{enable: false},
		},
		{
			structTag:       `flag:""`,
			expectedFlagTag: flagTag{enable: true},
		},
		{
			structTag: `flag:"name=foo type=string flat=true short=s env=true env-split=,"`,
			expectedFlagTag: flagTag{
				enable: true, name: "foo", flagType: "string", flat: true, shorthand: "s",
				enableEnv: true, envSplit: ",",
			},
		},
		{
			structTag:       `flag:"flat=true env=true"`,
			expectedFlagTag: flagTag{enable: true, flat: true, enableEnv: true},
		},
		{
			structTag:       `flag:"flat env"`,
			expectedFlagTag: flagTag{enable: true, flat: true, enableEnv: true},
		},
		{
			structTag:       `flag-usage:"foo,&\n"`,
			expectedFlagTag: flagTag{enable: true, usage: "foo,&\n"},
		},
		{
			structTag:       `flag:"" flag-usage:"foo"`,
			expectedFlagTag: flagTag{enable: true, usage: "foo"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.structTag, func(t *testing.T) {
			result := resolveFlagTag(reflect.StructTag(tc.structTag))

			assert.Equal(t, tc.expectedFlagTag, result)
		})
	}
}
