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

	if len(files) == 0 {
		files = append(files, "-")
	}

	for _, file := range files {
		if file == "-" {
			file = "/dev/stdin"
		}

		f, err := os.Stat(file)

		if err != nil {
			fmt.Fprintf(os.Stderr, "gocat: %s: No such file or directory\n", file)
			os.Exit(1)
		}

		if f.IsDir() {
			fmt.Fprintf(os.Stderr, "gocat: %s: Is a directory\n", file)
			os.Exit(1)
		}

		input, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer input.Close()

		gocat.CopyStream(input, os.Stdout, 1)
	}
}
