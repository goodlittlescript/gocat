package main

import (
	"flag"
	"fmt"
	"gocat"
	"os"
	"path/filepath"
)

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func main() {
	recursive := false
	flag.Usage = gocat.Usage(`
usage: gocp [options] source_file target_file

Copy files.

options:

`)
	help := flag.Bool("h", false, "print help")
	flag.BoolVar(&recursive, "R", recursive, "copy recursive")
	flag.BoolVar(&recursive, "r", recursive, "copy recursive")
	args := gocat.ParseToEnd()

	if *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.Usage()
		os.Exit(0)
	}

	if len(args) < 2 {
		flag.CommandLine.SetOutput(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	sources := args[:len(args)-1]
	target := args[len(args)-1]
	fileList, err := gocat.CopyList(sources, target, recursive)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gocp: %s\n", err)
		os.Exit(1)
	}

	for _, source_target := range fileList {
		sourcefile := source_target[0]
		targetfile := source_target[1]

		input, err := os.Open(sourcefile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer input.Close()

		os.MkdirAll(filepath.Dir(targetfile), os.ModePerm)
		output, err := os.Create(targetfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer output.Close()

		err = gocat.CopyStream(input, output, 1024)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
}
