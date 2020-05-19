# go-hybrid

![test & build](https://github.com/XSAM/go-hybrid/workflows/test%20&%20build/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/XSAM/go-hybrid/badge.svg?branch=master)](https://coveralls.io/github/XSAM/go-hybrid?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/XSAM/go-hybrid)](https://goreportcard.com/report/github.com/XSAM/go-hybrid)
[![Documentation](https://godoc.org/github.com/XSAM/go-hybrid?status.svg)](https://pkg.go.dev/mod/github.com/XSAM/go-hybrid)

This repository provides common utility packages for Golang. Such as `log`, `error` and `cmdutil` lib.

- [go-hybrid](#go-hybrid)
- [Installation](#installation)
- [Example](#example)
- [Packages](#packages)
  - [cmdutil](#cmdutil)
    - [Flag rules and variables](#flag-rules-and-variables)
  - [errorw](#errorw)
  - [log](#log)
  - [metadata](#metadata)
  - [builtinutil](#builtinutil)

# Installation

```bash
go get github.com/XSAM/go-hybrid
```

# [Example](_example/)

```
make
./bin/example
```

# Packages

## [cmdutil](https://pkg.go.dev/github.com/XSAM/go-hybrid/cmdutil)

`cmdutil` offer `ResolveFlagVariable` to register a struct field into [cobra](https://github.com/spf13/cobra) as a flag.

For example:

```golang
type Flag struct {
	Number      int           `flag:"env"`
	Duration    time.Duration `flag-usage:"change the duration"`
}

var flag Flag
var cmd cobra.Command

cmdutil.ResolveFlagVariable(&cmd, &flag)
```

`ResolveFlagVariable` will resolve struct through the struct tag. So there is no need to invoke `cmd.PersistentFlags()`, `ResolveFlagVariable` already do that for you.

Furthermore, `ResolveFlagVariable` extends the ability of `cobra`, it can automatically read the environment variable if you like. Simply add `env` to the `flag:""`.

The priority of the value assignment is:

`flag parameter > environment variable > default value`

For instance, setting a variable though a flag parameter and an environment variable at the same time, the variable value will be the flag parameter value.

### Flag rules and variables

Add `flag:""` or `flag-usage:""` to the struct tag and let `cmdutil` know that you want to resolve this variable.

Use `flag-usage` to add usage for a flag.

| key  | example             | description                                      |   |
|------|---------------------|--------------------------------------------------|---|
| env  | `env`, `env=true`   | read environment variable.                       |   |
| name | `name=foo`          | not like generated name? use it to overwrite it. |   |
| flat | `flat`, `flat=true` | ignore prefix name.                               |   |

Currently supported type: `bool`, `string`, `int`, `int32`, `int64`, `time.Duration`, `[]int`, `[]time.Duration`, `[]string`, `[]bool`, `map[string]string`, `map[string]int`

## [errorw](https://pkg.go.dev/github.com/XSAM/go-hybrid/errorw)

`errrow` is an error type which wraps fields and callstack. It implements `stackTracer`, `causer`, `GRPCStatus` and `zapcore.ObjectMarshaler` interface. 

Also, it uses gRPC status as an API error. So it can be directly returned as a gRPC error.

[grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway/blob/554b3dac4972c2957a8bc8e8ba15a241a6352b93/runtime/errors.go#L16) provides an approach that converting a gRPC error code into the corresponding HTTP response status. Therefore, you can use it even you want HTTP status codes.

## [log](https://pkg.go.dev/github.com/XSAM/go-hybrid/log)

`log` wraps [zap](go.uber.org/zap) as logger. It mainly provides `BgLogger()` and `Logger(ctx context.Context)` to access a concrete logger.

And, it provides customized preset config to control the log output style, such as `JSON` and `Text` style. You can use `environment` package to switch it. Check [this file](environment/service.go) for more details.

## [metadata](https://pkg.go.dev/github.com/XSAM/go-hybrid/metadata)

You can inject some const variables relevant to the program itself, such as *gitVersion*, *gitCommit*, *gitBranch* and *buildTime*. Then you can fetch these variables from `metadata.AppInfo`.

Also, you can use `cmdutil.Version` to add `version` command to `cobra`.

You can check out [Makefile](./Makefile) and learn how to inject these variables.

## [builtinutil](https://pkg.go.dev/github.com/XSAM/go-hybrid/builtinutil)

Providing `UserHomeDir`, `Recovery` and `WrappedGo`

`WrappedGo` wraps a goroutine with a recovery. So you will not worry about forget to recover a goroutine.