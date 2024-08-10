// Copyright (c) 2024  The Go-CoreUtils Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bep/golibsass/libsass"
	"github.com/urfave/cli/v2"

	clcli "github.com/go-corelibs/cli"
)

var (
	Version = "0.0.0"
	Release = "development"
)

func init() {

	cli.HelpFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "help",
		Aliases:  []string{"h", "usage"},
	}

	cli.VersionFlag = &cli.BoolFlag{
		Category: "General",
		Name:     "version",
		Aliases:  []string{"V"},
	}

	cli.FlagStringer = clcli.NewFlagStringer().
		PruneRepeats(true).
		PruneDefaultBools(true).
		DetailsOnNewLines(true).
		Make()

}

func main() {
	app := &cli.App{
		Name:        "gassc",
		Usage:       "Go sass compiler",
		Action:      action,
		Version:     Version + " (" + Release + ")",
		UsageText:   "gassc [options] <source.scss>",
		Description: "Simple libsass compiler.",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Category: "Outputs",
				Name:     "output-file",
				Usage:    "specify file to write, use \"-\" for stdout",
				Value:    "-",
				Aliases:  []string{"O"},
			},
			&cli.StringFlag{
				Category: "Outputs",
				Name:     "output-style",
				Usage:    "nested, expanded, compact or compressed",
				Value:    "nested",
				Aliases:  []string{"S"},
			},
			&cli.BoolFlag{
				Category: "Outputs",
				Name:     "no-source-map",
				Usage:    "do not include source-map output",
				Aliases:  []string{"M"},
			},
			&cli.StringSliceFlag{
				Category: "Settings",
				Name:     "include-path",
				Usage:    "add one (or more) include paths",
				Aliases:  []string{"I"},
			},
			&cli.BoolFlag{
				Category: "Settings",
				Name:     "sass-syntax",
				Usage:    "use sass instead of scss syntax",
				Aliases:  []string{"A"},
			},
			&cli.IntFlag{
				Category: "Settings",
				Name:     "precision",
				Usage:    "floating point precision",
				Value:    10,
				Aliases:  []string{"P"},
			},
			&cli.BoolFlag{
				Category: "Settings",
				Name:     "release",
				Usage:    "same as: -M -S=compressed",
			},
		},
		HideHelpCommand:      true,
		EnableBashCompletion: true,
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func action(ctx *cli.Context) (err error) {
	if ctx.NArg() != 1 {
		cli.ShowAppHelpAndExit(ctx, 1)
	}
	if err = process(ctx, ctx.Args().First()); err != nil {
		return
	}
	return
}

func process(ctx *cli.Context, src string) (err error) {

	var data []byte
	if data, err = os.ReadFile(src); err != nil {
		return
	}
	contents := string(data)

	releaseMode := ctx.Bool("release")

	outputFile := ctx.Path("output-file")
	if outputFile == "" {
		outputFile = "-"
	}
	outputStyle := ctx.String("output-style")
	if releaseMode {
		outputStyle = "compressed"
	}

	includePaths := ctx.StringSlice("include-path")
	includePaths = append(includePaths, filepath.Dir(src))
	options := libsass.Options{
		IncludePaths: includePaths,
		Precision:    ctx.Int("precision"),
		SassSyntax:   ctx.Bool("sass-syntax"),
		OutputStyle:  libsass.ParseOutputStyle(outputStyle),
	}

	if !releaseMode && !ctx.Bool("no-source-map") {
		if outputFile == "-" {
			options.SourceMapOptions = libsass.SourceMapOptions{
				Contents:       true,
				OmitURL:        true,
				EnableEmbedded: true,
			}
		} else {
			options.SourceMapOptions = libsass.SourceMapOptions{
				Filename:       outputFile + ".map",
				Contents:       true,
				OmitURL:        true,
				EnableEmbedded: false,
			}
		}
	}

	var transpiler libsass.Transpiler
	if transpiler, err = libsass.New(options); err != nil {
		err = fmt.Errorf("error constructing transpilier: %v", err)
		return
	}

	var result libsass.Result
	if result, err = transpiler.Execute(contents); err != nil {
		err = fmt.Errorf("error transipiling: %v", err)
		return
	}

	if outputFile == "" || outputFile == "-" {
		_, _ = fmt.Fprint(os.Stdout, result.CSS)
		return
	}
	if err = os.WriteFile(outputFile, []byte(result.CSS), 0660); err != nil {
		err = fmt.Errorf("error writing output-file: %v", err)
		return
	}
	if result.SourceMapFilename != "" && result.SourceMapContent != "" {
		if err = os.WriteFile(result.SourceMapFilename, []byte(result.SourceMapContent), 0660); err != nil {
			err = fmt.Errorf("error writing sourcemap file: %v - %v", result.SourceMapFilename, err)
			return
		}
	}
	return
}
