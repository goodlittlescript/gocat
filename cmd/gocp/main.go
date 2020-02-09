package main

import (
	"flag"
	"fmt"
	"gocat"
	"os"
)

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func main() {
	recursive := false
	flag.Usage = gocat.Usage(`
usage: gocp [options] source_file target_file
       gocp [options] source_file ... target

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
		fmt.Fprintf(os.Stderr, "%s\n", gocat.Desc())
		os.Exit(1)
	}

	num_fail := 0
	err := gocat.CopyFiles(args, recursive, gocat.CopyStream, func(failure error) {
		fmt.Fprintf(os.Stderr, "%s\n", failure)
		num_fail += 1
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}

	if num_fail > 0 {
		os.Exit(1)
	}
}
