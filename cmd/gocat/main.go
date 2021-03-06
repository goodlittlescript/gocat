package main

import (
	"flag"
	"github.com/goodlittlescript/gocat"
	"os"
)

func main() {
	flag.Usage = gocat.Usage(`
usage: gocat [options] [file ...]

Concatenate and print files.

options:

`)
	help := flag.Bool("h", false, "print help")
	flag.Bool("u", false, "unbuffer output") // appears to behave as true, always
	files := gocat.ParseToEnd()

	if *help {
		flag.CommandLine.SetOutput(os.Stdout)
		flag.Usage()
		os.Exit(0)
	}

	err := gocat.CatFiles(files, os.Stdout, gocat.CopyStream)
	gocat.Check(err)
}
