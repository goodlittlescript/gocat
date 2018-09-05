package main

import "fmt"
import "os"
import "io"
import "github.com/spf13/pflag"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var help bool
	var unbuffer bool
	var desc = `
usage: gocat [options] FILES...

Concatenate and print files.

options:

`
	pflag.BoolVarP(&help, "help", "h", false, "print help")
	pflag.BoolVarP(&unbuffer, "unbuffer", "u", false, "unbuffer output")

	pflag.Usage = func() {
		fmt.Printf(desc[1:])
		pflag.CommandLine.SetOutput(os.Stdout)
		pflag.PrintDefaults()
		fmt.Println()
	}
	pflag.Parse()

	if help {
		pflag.Usage()
		os.Exit(0)
	}

	files := pflag.Args()
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
			check(err)
			defer input.Close()
		}

		buffer := make([]byte, 100)
		for {
			nbytes, err := input.Read(buffer)
			if err != nil {
				if err != io.EOF {
					check(err)
				}

				break
			}

			_, err = os.Stdout.Write(buffer[:nbytes])
			if err != nil {
				check(err)
			}
		}
	}
}
