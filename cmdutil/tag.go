package cmdutil

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	tagKeySeparator = "="
)

type flagTag struct {
	// Enable to generate flag
	enable bool

	name *string
	// Without prefix name
	flagType  string
	flat      bool
	required  bool
	shorthand string
	usage     string

	enableEnv bool
	envSplit  string
}

func resolveFlagTag(structTag reflect.StructTag) flagTag {
	tag, ok := structTag.Lookup("flag")
	if !ok && structTag.Get("flag-usage") == "" {
		return flagTag{enable: false}
	}

	flagKV := make(map[string]string)

	f := func(c rune) bool {
		return unicode.IsSpace(c)
	}

	// Splitting string by space but considering quoted section
	items := strings.FieldsFunc(tag, f)
	// Create and fill the map
	for _, item := range items {
		vals := strings.SplitN(item, tagKeySeparator, 2)

		var key, value string
		switch len(vals) {
		case 1:
			key = vals[0]
		case 2:
			key = vals[0]
			value = vals[1]
		}
		if key != "" {
			flagKV[key] = value
		}
	}

	// Fill struct
	var flat, enableEnv, required bool
	if v, ok := flagKV["flat"]; ok {
		flat = parseBool(v)
	}
	if v, ok := flagKV["env"]; ok {
		enableEnv = parseBool(v)
	}
	if v, ok := flagKV["required"]; ok {
		required = parseBool(v)
	}
	var name *string
	if v, ok := flagKV["name"]; ok {
		name = newString(v)
	}

	return flagTag{
		enable:    true,
		name:      name,
		flat:      flat,
		required:  required,
		flagType:  flagKV["type"],
		shorthand: flagKV["short"],
		// Prevent parse error since usage may have ','
		usage:     structTag.Get("flag-usage"),
		enableEnv: enableEnv,
		envSplit:  flagKV["env-split"],
	}
}

func parseBool(v string) bool {
	if v == "" {
		return true
	}
	result, _ := strconv.ParseBool(v)
	return result
}

func newString(b string) *string {
	return &b
}
