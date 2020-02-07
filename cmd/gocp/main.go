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

	if len(args) < 2  {
		flag.CommandLine.SetOutput(os.Stderr)
		flag.Usage()
		os.Exit(1)
	}

	fileList := [][]string{}

	if recursive {
		sourceDir := args[0]
		targetDir := args[1]
		err := filepath.Walk(sourceDir, func(sourcefile string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() {
				   return nil
				}

				rel, err := filepath.Rel(sourceDir, sourcefile)
				if err != nil {
					return err
				}

				targetfile := filepath.Join(targetDir, rel)
				fileList = append(fileList, []string{sourcefile, targetfile})

				return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	} else {
		sourcefile := args[0]
		targetfile := args[1]
		fileList = append(fileList, []string{sourcefile, targetfile})
	}

	for _, source_target := range fileList {
		sourcefile := source_target[0]
		targetfile := source_target[1]

		source, err := os.Stat(sourcefile)

		if err != nil {
			fmt.Fprintf(os.Stderr, "gocp: %s: No such file or directory\n", sourcefile)
			os.Exit(1)
		}

		if source.IsDir() {
			fmt.Fprintf(os.Stderr, "gocp: %s is a directory (not copied).\n", sourcefile)
			os.Exit(1)
		}

		target, err := os.Stat(targetfile)

		if err == nil {
			if target.IsDir() {
				targetfile = filepath.Join(targetfile)
			}
		}

		if sourcefile == targetfile {
			fmt.Fprintf(os.Stderr, "gocp: %s and %s are identical (not copied).\n", sourcefile, targetfile)
			os.Exit(1)
		}
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
