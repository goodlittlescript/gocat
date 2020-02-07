package main

import (
	"flag"
	"fmt"
	"gocat"
	"os"
)

func main() {
	flag.Usage = gocat.Usage(`
usage: gocat [options] [file ...]

Concatenate and print files.

options:

`)
	help := flag.Bool("h", false, "print help")
	flag.Bool("u", false, "unbuffer output")
	files := gocat.ParseToEnd()

	if *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.Usage()
		os.Exit(0)
	}

	err := gocat.CatFiles(files, os.Stdout, gocat.CopyStream)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
