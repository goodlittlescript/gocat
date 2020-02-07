package gocat

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Desc() string {
	curr := flag.CommandLine.Output()
	buf := bytes.NewBufferString("")
	flag.CommandLine.SetOutput(buf)
	flag.Usage()
	flag.CommandLine.SetOutput(curr)
	return strings.Split(buf.String(), "\n")[0]
}

func Usage(desc string) func() {
	return func() {
		output := flag.CommandLine.Output()
		fmt.Fprintf(output, desc[1:])
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(output, "  -%-6s %s\n", f.Name, f.Usage)
		})
		fmt.Fprintf(output, "\n")
	}
}

func Parse() []string {
	flag.Parse()
	return flag.Args()
}

func ParseToEnd() []string {
	flag.Parse()
	args := make([]string, 0)
	for i := len(os.Args) - len(flag.Args()) + 1; i < len(os.Args); {
		if i > 1 && os.Args[i-2] == "--" {
			break
		}
		args = append(args, flag.Arg(0))
		flag.CommandLine.Parse(os.Args[i:])
		i += 1 + len(os.Args[i:]) - len(flag.Args())
	}
	args = append(args, flag.Args()...)
	return args
}

func CopyStream(input io.Reader, output io.Writer) error {
	buffer := make([]byte, 1024)
	for {
		n, rerr := input.Read(buffer)

		_, werr := output.Write(buffer[:n])
		if werr != nil {
			return werr
		}

		if rerr != nil {
			if rerr == io.EOF {
				return nil
			}
			return rerr
		}
	}
}

func CatFiles(files []string, output io.Writer, copyfunc func(io.Reader, io.Writer) error) error {
	if len(files) == 0 {
		files = append(files, "-")
	}

	for i, file := range files {
		if file == "-" {
			file = "/dev/stdin"
			files[i] = file
		}

		f, err := os.Stat(file)

		if err != nil {
			return fmt.Errorf("%s: No such file or directory", file)
		}

		if f.IsDir() {
			return fmt.Errorf("%s: Is a directory", file)
		}
	}

	for _, file := range files {
		input, err := os.Open(file)
		if err != nil {
			return err
		}
		defer input.Close()

		err = copyfunc(input, output)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyFiles(sources []string, target string, recursive bool, copyfunc func(io.Reader, io.Writer) error) error {
	fileList := [][2]string{}

	// Assemble list of {sourceFile, sourceDir}
	// Keep the dir so that relative paths may be calculated if needed.
	for _, source := range sources {
		if recursive {
			sourceDir := source
			err := filepath.Walk(sourceDir, func(sourceFile string, f os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if f.IsDir() {
					return nil
				}

				fileList = append(fileList, [2]string{sourceFile, sourceDir})
				return nil
			})

			if err != nil {
				return err
			}
		} else {
			sourceFile := source
			fileList = append(fileList, [2]string{sourceFile, filepath.Dir(source)})
		}
	}

	// Convert to {sourceFile, targetFile}
	if len(sources) > 1 || recursive {
		if !recursive {
			t, err := os.Stat(target)

			if err != nil || !t.IsDir() {
				return fmt.Errorf("target '%s' is not a directory", target)
			}
		}

		for i, item := range fileList {
			sourceFile := item[0]
			sourceDir := item[1]

			rel, err := filepath.Rel(sourceDir, sourceFile)
			if err != nil {
				return err
			}

			targetFile := filepath.Join(target, rel)
			fileList[i][1] = targetFile
		}
	} else {
		fileList[0][1] = target
	}

	// Check each to confirm copy will work
	for _, item := range fileList {
		sourceFile := item[0]
		targetFile := item[1]

		s, err := os.Stat(sourceFile)

		if err != nil {
			return fmt.Errorf("%s: No such file or directory", sourceFile)
		}

		if s.IsDir() {
			return fmt.Errorf("%s is a directory (not copied).", sourceFile)
		}

		t, err := os.Stat(targetFile)

		if err == nil {
			if t.IsDir() {
				targetFile = filepath.Join(targetFile)
			}
		}

		if sourceFile == targetFile {
			return fmt.Errorf("%s and %s are identical (not copied).", sourceFile, targetFile)
		}
	}

	// Copy files
	for _, item := range fileList {
		sourceFile := item[0]
		targetFile := item[1]

		input, err := os.Open(sourceFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		defer input.Close()

		os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)
		output, err := os.Create(targetFile)
		if err != nil {
			return err
		}
		defer output.Close()

		err = copyfunc(input, output)
		if err != nil {
			return err
		}
	}

	return nil
}
