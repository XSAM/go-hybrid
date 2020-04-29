package cmdutil

import (
	"os"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type testFlag struct {
	Ignore  string
	Ignore2 string `json:"ignore_2"`
	// Cannot return value obtained from unexported field or method
	ignore3 string `flag:""`
	// Don't have flag tag
	Ignore4 TestIgnore

	String string            `flag:"flat=true short=s env=true env-split=," flag-usage:"foo"`
	Array  []string          `flag:""`
	Map    map[string]string `flag:""`

	OverwriteName   string     `flag:"name=new-name"`
	OverwritePrefix testPrefix `flag:"name=new-prefix"`

	Short      string `flag:""`
	ShortEnv   string `flag:"env"`
	ShortFlat  string `flag:"flat"`
	ShortUsage string `flag-usage:"foo"`

	TestInline `flag:""`
	Flat       testFlat `flag:""`
	// Set prefix is empty
	Flat2 testFlag2 `flag:"name="`
	Type  testType  `flag:""`
}

type testFlat struct {
	Flat string `flag:"flat"`
}

type testFlag2 struct {
	Flat2 string `flag:""`
}

type TestIgnore struct {
	Ignore string `flag:""`
}

type testType struct {
	Int            int               `flag:""`
	String         string            `flag:""`
	Bool           bool              `flag:""`
	IntSlice       []int             `flag:""`
	StringSlice    []string          `flag:""`
	BoolSlice      []bool            `flag:""`
	StringToString map[string]string `flag:""`
	StringToInt    map[string]int    `flag:""`
}

type testPrefix struct {
	Prefix string `flag:""`
}

type TestInline struct {
	Inline string `flag:""`
}

func Test_resolveFlags(t *testing.T) {
	f := testFlag{
		String:  "normal",
		Array:   []string{"foo"},
		Map:     map[string]string{"foo": "bar"},
		ignore3: "",
	}

	var result flags
	result = resolveFlags(&f, result, "", 0)
	assert.Equal(t, flags{
		{Name: "string", FullName: "string", FullEnv: "STRING", EnableEnv: true, Shorthand: "s", Usage: "foo", EnvSplit: ",", Type: "string", Value: "normal", Pointer: &f.String},
		{Name: "array", FullName: "array", FullEnv: "ARRAY", Type: "string-slice", Value: []string{"foo"}, Pointer: &f.Array},
		{Name: "map", FullName: "map", FullEnv: "MAP", Type: "string-to-string", Value: map[string]string{"foo": "bar"}, Pointer: &f.Map},

		{Name: "new-name", FullName: "new-name", FullEnv: "NEW_NAME", Type: "string", Value: "", Pointer: &f.OverwriteName},
		{Name: "prefix", FullName: "new-prefix-prefix", FullEnv: "NEW_PREFIX_PREFIX", Type: "string", Value: "", Pointer: &f.OverwritePrefix.Prefix},

		{Name: "short", FullName: "short", FullEnv: "SHORT", Type: "string", Value: "", Pointer: &f.Short},
		{Name: "short-env", FullName: "short-env", FullEnv: "SHORT_ENV", EnableEnv: true, Type: "string", Value: "", Pointer: &f.ShortEnv},
		{Name: "short-flat", FullName: "short-flat", FullEnv: "SHORT_FLAT", Type: "string", Value: "", Pointer: &f.ShortFlat},
		{Name: "short-usage", FullName: "short-usage", FullEnv: "SHORT_USAGE", Type: "string", Value: "", Usage: "foo", Pointer: &f.ShortUsage},

		{Name: "inline", FullName: "test-inline-inline", FullEnv: "TEST_INLINE_INLINE", Type: "string", Value: "", Pointer: &f.TestInline.Inline},

		{Name: "flat", FullName: "flat", FullEnv: "FLAT", Type: "string", Value: "", Pointer: &f.Flat.Flat},

		{Name: "flat2", FullName: "flat2", FullEnv: "FLAT2", Type: "string", Value: "", Pointer: &f.Flat2.Flat2},

		{Name: "int", FullName: "type-int", FullEnv: "TYPE_INT", Type: "int", Value: 0, Pointer: &f.Type.Int},
		{Name: "string", FullName: "type-string", FullEnv: "TYPE_STRING", Type: "string", Value: "", Pointer: &f.Type.String},
		{Name: "bool", FullName: "type-bool", FullEnv: "TYPE_BOOL", Type: "bool", Value: false, Pointer: &f.Type.Bool},
		{Name: "int-slice", FullName: "type-int-slice", FullEnv: "TYPE_INT_SLICE", Type: "int-slice", Value: ([]int)(nil), Pointer: &f.Type.IntSlice},
		{Name: "string-slice", FullName: "type-string-slice", FullEnv: "TYPE_STRING_SLICE", Type: "string-slice", Value: ([]string)(nil), Pointer: &f.Type.StringSlice},
		{Name: "bool-slice", FullName: "type-bool-slice", FullEnv: "TYPE_BOOL_SLICE", Type: "bool-slice", Value: ([]bool)(nil), Pointer: &f.Type.BoolSlice},
		{Name: "string-to-string", FullName: "type-string-to-string", FullEnv: "TYPE_STRING_TO_STRING", Type: "string-to-string", Value: (map[string]string)(nil), Pointer: &f.Type.StringToString},
		{Name: "string-to-int", FullName: "type-string-to-int", FullEnv: "TYPE_STRING_TO_INT", Type: "string-to-int", Value: (map[string]int)(nil), Pointer: &f.Type.StringToInt},
	}, result)

	// Can be registered
	err := ResolveFlagVariable(&cobra.Command{}, &f)
	assert.NoError(t, err)
}

func Test_toSnake(t *testing.T) {
	testCases := []struct {
		name     string
		expected string
	}{
		{
			name:     "InfluxDB",
			expected: "influx-db",
		},
		{
			name:     "InfluxDBV2",
			expected: "influx-dbv2",
		},
		{
			name:     "fooBarVer",
			expected: "foo-bar-ver",
		},
		{
			name:     "EtcdEndpoints",
			expected: "etcd-endpoints",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := toSnake(tc.name)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_resolveCobraType(t *testing.T) {
	testCases := []struct {
		name         string
		field        interface{}
		tag          flagTag
		expectedType string
	}{
		{
			name: "overwrite field type",
			field: struct {
				M string
			}{},
			tag:          flagTag{flagType: "new"},
			expectedType: "new",
		},
		{
			name: "int",
			field: struct {
				M int
			}{},
			expectedType: "int",
		},
		{
			name: "bool",
			field: struct {
				M bool
			}{},
			expectedType: "bool",
		},
		{
			name: "string",
			field: struct {
				M string
			}{},
			expectedType: "string",
		},
		{
			name: "int-slice",
			field: struct {
				M []int
			}{},
			expectedType: "int-slice",
		},
		{
			name: "bool-slice",
			field: struct {
				M []bool
			}{},
			expectedType: "bool-slice",
		},
		{
			name: "string-slice",
			field: struct {
				M []string
			}{},
			expectedType: "string-slice",
		},
		{
			name: "map int bool",
			field: struct {
				M map[int]bool
			}{},
			expectedType: "int-to-bool",
		},
		{
			name: "map int string",
			field: struct {
				M map[int]string
			}{},
			expectedType: "int-to-string",
		},
		{
			name: "map bool bool",
			field: struct {
				M map[bool]bool
			}{},
			expectedType: "bool-to-bool",
		},
		{
			name: "map string int",
			field: struct {
				M map[string]int
			}{},
			expectedType: "string-to-int",
		},
		{
			name: "invalid: complex type map",
			field: struct {
				M map[bool]interface{}
			}{},
			expectedType: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := reflect.TypeOf(tc.field).Field(0)

			result := resolveCobraType(field, tc.tag)
			assert.Equal(t, tc.expectedType, result)
		})
	}
}

func TestResolveFlagVariableWithInvalid(t *testing.T) {
	var err error
	// Require pointer type
	m := struct {
		Name string
	}{}
	err = ResolveFlagVariable(nil, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "flag variable require pointer type")

	// Not supported type
	n := struct {
		M map[int]string `flag:""`
	}{}
	err = ResolveFlagVariable(nil, &n)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not supported flag type: int-to-string")
}

func TestResolveFlagVariableWithEnv(t *testing.T) {
	// Monkey patch
	monkey.Patch(os.Getenv, func(key string) string {
		switch key {
		case "NAME":
			return "value"
		case "SLICE":
			return "foo,bar"
		case "SLICE2":
			return "foo;bar"
		case "SLICE3":
			return "3,4"
		}
		return ""
	})
	defer monkey.Unpatch(os.Getenv)

	m := struct {
		Name  string   `flag:"env"`
		Slice []string `flag:"env"`
		// Special value separator
		Slice2 []string `flag:"env env-split=;"`
		// Overwrite default values
		Slice3 []int `flag:"env"`
	}{
		Slice3: []int{1, 2},
	}

	err := ResolveFlagVariable(&cobra.Command{}, &m)

	assert.NoError(t, err)
	assert.Equal(t, "value", m.Name)
	assert.Equal(t, []string{"foo", "bar"}, m.Slice)
	assert.Equal(t, []string{"foo", "bar"}, m.Slice2)
	assert.Equal(t, []int{3, 4}, m.Slice3)
}

func TestResolveFlagVariableWithEnvAndUsage(t *testing.T) {
	// Monkey patch
	monkey.Patch(os.Getenv, func(key string) string {
		switch key {
		case "NAME":
			return "value2"
		case "NAME2":
			return "value3"
		}
		return ""
	})
	defer monkey.Unpatch(os.Getenv)

	m := struct {
		Name  string `flag:"env"`
		Name2 string `flag:"env" flag-usage:"usage"`
	}{}

	cmd := cobra.Command{}
	err := ResolveFlagVariable(&cmd, &m)

	assert.NoError(t, err)

	assert.Equal(t, "[env NAME]", cmd.Flag("name").Usage)
	assert.Equal(t, "value2", m.Name)

	assert.Equal(t, "usage [env NAME2]", cmd.Flag("name2").Usage)
	assert.Equal(t, "value3", m.Name2)
}

func TestResolveFlagVariableWithDefaultValue(t *testing.T) {
	m := struct {
		Name string `flag:""`
	}{
		Name: "foo",
	}

	cmd := cobra.Command{}
	err := ResolveFlagVariable(&cmd, &m)

	assert.NoError(t, err)
	assert.Equal(t, "foo", cmd.Flag("name").Value.String())
}

func TestResolveFlagVariableWithConflictedName(t *testing.T) {
	m := struct {
		Name  string `flag:""`
		Name2 string `flag:"name=name"`
	}{}

	cmd := cobra.Command{}
	err := ResolveFlagVariable(&cmd, &m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicated flag full name")
}

func Test_resolveFlagsDepth(t *testing.T) {
	FlagMaxDepth = 2
	defer func() {
		FlagMaxDepth = defaultFlagMaxDepth
	}()

	m := struct {
		Name string `flag:""`
		M    struct {
			Name string `flag:""`
			M    struct {
				// Should be ignore
				Name string `flag:""`
			} `flag:""`
		} `flag:""`
	}{}

	cmd := cobra.Command{}
	err := ResolveFlagVariable(&cmd, &m)
	assert.NoError(t, err)

	var flags []*pflag.Flag
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flags = append(flags, flag)
	})

	assert.Len(t, flags, 2)
}

func TestResolveFlagVariableWithWrongEnv(t *testing.T) {
	// Monkey patch
	monkey.Patch(os.Getenv, func(key string) string {
		switch key {
		case "INT":
			return "value"
		}
		return ""
	})
	defer monkey.Unpatch(os.Getenv)

	m := struct {
		Int int `flag:"env"`
	}{}

	err := ResolveFlagVariable(&cobra.Command{}, &m)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set env value")
}

func TestResolveFlagVariableWithWrongEnvSplit(t *testing.T) {
	// Monkey patch
	monkey.Patch(os.Getenv, func(key string) string {
		switch key {
		case "SLICE":
			return "value;value"
		}
		return ""
	})
	defer monkey.Unpatch(os.Getenv)

	m := struct {
		Slice []int `flag:"env env-split=;"`
	}{}

	err := ResolveFlagVariable(&cobra.Command{}, &m)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "set env value with delimiter")
}
