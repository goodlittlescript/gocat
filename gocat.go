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

func CopyStream(input io.Reader, output io.Writer, chunk_size int) error {
	buffer := make([]byte, chunk_size)
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

func CopyList(sources []string, target string, recursive bool) ([][2]string, error) {
	fileList := [][2]string{}

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
				return nil, err
			}
		} else {
			sourceFile := source
			fileList = append(fileList, [2]string{sourceFile, filepath.Dir(source)})
		}
	}

	if len(sources) > 1 || recursive {
		if !recursive {
			t, err := os.Stat(target)

			if err != nil || !t.IsDir() {
				return nil, fmt.Errorf("target '%s' is not a directory\n", target)
			}
		}

		for i, item := range fileList {
			sourceFile := item[0]
			sourceDir := item[1]

			rel, err := filepath.Rel(sourceDir, sourceFile)
			if err != nil {
				return nil, err
			}

			targetFile := filepath.Join(target, rel)
			fileList[i][1] = targetFile
		}
	} else {
		fileList[0][1] = target
	}

	for _, item := range fileList {
		sourceFile := item[0]
		targetFile := item[1]

		s, err := os.Stat(sourceFile)

		if err != nil {
			return nil, fmt.Errorf("%s: No such file or directory\n", sourceFile)
		}

		if s.IsDir() {
			return nil, fmt.Errorf("%s is a directory (not copied).\n", sourceFile)
		}

		t, err := os.Stat(targetFile)

		if err == nil {
			if t.IsDir() {
				targetFile = filepath.Join(targetFile)
			}
		}

		if sourceFile == targetFile {
			return nil, fmt.Errorf("%s and %s are identical (not copied).\n", sourceFile, targetFile)
		}
	}

	return fileList, nil
}
