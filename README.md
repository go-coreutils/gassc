[![Go](https://img.shields.io/badge/Go-v1.22.6-blue.svg)](https://go.dev)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/go-coreutils/gassc)
[![GoReportCard](https://goreportcard.com/badge/github.com/go-coreutils/gassc)](https://goreportcard.com/report/github.com/go-coreutils/gassc)

# gassc - Go sass compiler

This is the SCSS compiler used by [Go-Enjin] theme features.

# Installation

``` shell
# note: requires Go 1.22.6
$ go install github.com/go-coreutils/gassc@latest
```

# Usage

``` shell
> gassc --help
NAME:
   gassc - Go sass compiler

USAGE:
   gassc [options] <source.scss>

VERSION:
   0.3.0 (trunk)

DESCRIPTION:
   Simple libsass compiler.

GLOBAL OPTIONS:
   General

   --help, -h, --usage
   --version, -V

   Outputs

   --no-source-map, -M             do not include source-map output
   --output-file value, -O value   specify file to write, use "-" for stdout
                                     (default: "-")
   --output-style value, -S value  nested, expanded, compact or compressed
                                     (default: "nested")

   Settings

   --include-path value, -I value  add one (or more) include paths
   --precision value, -P value     floating point precision
                                     (default: 10)
   --release                       same as: -M -S=compressed
   --sass-syntax, -A               use sass instead of scss syntax
```

# License

```
Copyright 2024  The Go-CoreUtils Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use file except in compliance with the License.
You may obtain a copy of the license at

 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

[Go-Enjin]: https://go-enjin.org
