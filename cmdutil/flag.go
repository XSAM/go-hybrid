package cmdutil

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/XSAM/go-hybrid/errorw"
)

type flags []flag

type flag struct {
	// e.g. foo
	Name string
	// With struct hierarchy prefix
	// e.g. prefix-foo
	FullName  string
	Shorthand string
	Usage     string

	// Enable env
	EnableEnv bool
	// With struct hierarchy prefix
	// e.g. PREFIX_FOO
	FullEnv  string
	EnvSplit string

	Type    string
	Value   interface{}
	Pointer interface{}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnake(name string) string {
	snake := matchFirstCap.ReplaceAllString(name, "${1}-${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}-${2}")
	return strings.ToLower(snake)
}

func resolveCobraType(field reflect.StructField, tag flagTag) string {
	if tag.flagType != "" {
		return tag.flagType
	}

	switch field.Type.Kind() {
	case reflect.Int:
		return "int"
	case reflect.Bool:
		return "bool"
	case reflect.String:
		return "string"
	case reflect.Slice, reflect.Array:
		switch field.Type.Elem().Kind() {
		case reflect.Int:
			return "int-slice"
		case reflect.Bool:
			return "bool-slice"
		case reflect.String:
			return "string-slice"
		}
	case reflect.Map:
		var keyType, valueType string
		switch field.Type.Key().Kind() {
		case reflect.Int:
			keyType = "int"
		case reflect.Bool:
			keyType = "bool"
		case reflect.String:
			keyType = "string"
		}

		switch field.Type.Elem().Kind() {
		case reflect.Int:
			valueType = "int"
		case reflect.Bool:
			valueType = "bool"
		case reflect.String:
			valueType = "string"
		}

		if keyType != "" && valueType != "" {
			return fmt.Sprintf("%s-to-%s", keyType, valueType)
		}
	}
	return ""
}

func resolveFlags(obj interface{}, flags flags, namePrefix string, depth int) flags {
	if depth > FlagMaxDepth {
		return flags
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		v = v.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Prevent panic "cannot return value obtained from unexported field or method"
		if !v.Field(i).Addr().CanInterface() {
			continue
		}

		tag := resolveFlagTag(field.Tag)
		if !tag.enable {
			continue
		}
		name := resolveFieldName(field, tag)
		fullName := genFullName(tag.flat, namePrefix, name)

		switch field.Type.Kind() {
		case reflect.Struct:
			flags = resolveFlags(v.Field(i).Addr().Interface(), flags, fullName, depth+1)
		default:
			flagType := resolveCobraType(field, tag)
			flag := flag{
				Name:      name,
				FullName:  fullName,
				Shorthand: tag.shorthand,
				Usage:     tag.usage,
				EnableEnv: tag.enableEnv,
				FullEnv:   genEnv(fullName),
				EnvSplit:  tag.envSplit,
				Type:      flagType,
				Value:     v.Field(i).Interface(),
				Pointer:   v.Field(i).Addr().Interface(),
			}

			flags = append(flags, flag)
		}
	}
	return flags
}

func resolveFieldName(field reflect.StructField, tag flagTag) string {
	if tag.name != nil {
		return *tag.name
	}

	// Use field name
	return toSnake(field.Name)
}

func genFullName(flat bool, namePrefix, name string) string {
	if flat {
		return name
	}

	var fullName string
	if namePrefix != "" {
		fullName = namePrefix + "-" + name
	} else {
		fullName = name
	}
	return fullName
}

// genEnv replace "-" to "_", and upper all character
func genEnv(name string) string {
	return strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}

// ResolveFlagVariable register persistent flags and env via tags in struct
func ResolveFlagVariable(cmd *cobra.Command, f interface{}) (err error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Ptr {
		return errorw.NewMessage("flag variable require pointer type")
	}

	var flags flags
	flags = resolveFlags(f, flags, "", 1)

	// Check full name conflict
	set := make(map[string]struct{})
	for _, v := range flags {
		if _, ok := set[v.FullName]; ok {
			return errorw.NewMessagef("duplicated flag full name: %s", v.FullName)
		} else {
			set[v.FullName] = struct{}{}
		}
	}

	// Register flags to cobra
	for _, v := range flags {
		switch v.Type {
		case "bool":
			cmd.PersistentFlags().BoolVarP(v.Pointer.(*bool), v.FullName, v.Shorthand, v.Value.(bool), v.Usage)
		case "string":
			cmd.PersistentFlags().StringVarP(v.Pointer.(*string), v.FullName, v.Shorthand, v.Value.(string), v.Usage)
		case "int":
			cmd.PersistentFlags().IntVarP(v.Pointer.(*int), v.FullName, v.Shorthand, v.Value.(int), v.Usage)
		case "int-slice":
			cmd.PersistentFlags().IntSliceVarP(v.Pointer.(*[]int), v.FullName, v.Shorthand, v.Value.([]int), v.Usage)
		case "string-slice":
			cmd.PersistentFlags().StringSliceVarP(v.Pointer.(*[]string), v.FullName, v.Shorthand, v.Value.([]string), v.Usage)
		case "bool-slice":
			cmd.PersistentFlags().BoolSliceVarP(v.Pointer.(*[]bool), v.FullName, v.Shorthand, v.Value.([]bool), v.Usage)
		case "string-to-string":
			cmd.PersistentFlags().StringToStringVarP(v.Pointer.(*map[string]string), v.FullName, v.Shorthand, v.Value.(map[string]string), v.Usage)
		case "string-to-int":
			cmd.PersistentFlags().StringToIntVarP(v.Pointer.(*map[string]int), v.FullName, v.Shorthand, v.Value.(map[string]int), v.Usage)
		default:
			return errorw.NewMessagef("not supported flag type: %s", v.Type)
		}
	}

	// Register env
	err = registerEnv(cmd, flags)
	if err != nil {
		return errorw.Wrap(err, "register env value")
	}
	return nil
}

// Register env
func registerEnv(cmd *cobra.Command, flags flags) error {
	for _, v := range flags {
		if !v.EnableEnv {
			continue
		}

		f := cmd.Flag(v.FullName)

		if f.Usage == "" {
			f.Usage = fmt.Sprintf("[env %v]", v.FullEnv)
		} else {
			f.Usage = fmt.Sprintf("%v [env %v]", f.Usage, v.FullEnv)
		}
		if value := os.Getenv(v.FullEnv); value != "" {
			// Customized env separator
			if v.EnvSplit != "" {
				strings.Split(value, v.EnvSplit)

				for _, value := range strings.Split(value, v.EnvSplit) {
					err := f.Value.Set(value)
					if err != nil {
						return errorw.Wrap(err, "set env value with delimiter").
							WithField("name", v.FullName).
							WithField("value", value).
							WithField("delimiter", v.EnvSplit)
					}
				}
			} else {
				err := f.Value.Set(value)
				if err != nil {
					return errorw.Wrap(err, "set env value").
						WithField("name", v.FullName).
						WithField("value", value)
				}
			}
		}
	}
	return nil
}
