package main

import (
	"flag"
	"fmt"
	"gocat"
	"os"
)

func main() {
	flag.Usage = gocat.Usage(`
usage: gocat [options] FILES...

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

	if len(files) == 0 {
		files = append(files, "-")
	}

	var input *os.File
	var err error
	for _, file := range files {
		if file == "-" {
			input = os.Stdin
		} else {
			input, err = os.Open(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer input.Close()
		}

		gocat.CopyStream(input, os.Stdout, 1)
	}
}
