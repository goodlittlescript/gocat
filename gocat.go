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

type CopyFunc func(io.Reader, io.Writer) error

func CopyFiles(args []string, recursive bool, copyfunc CopyFunc) error {
	if len(args) == 0 {
		return fmt.Errorf("missing file operand")
	}

	if len(args) == 1 {
		return fmt.Errorf("missing destination file operand after '%s'", args[0])
	}

	sources := args[:len(args)-1]
	target := args[len(args)-1]

	// 1) len(args) = 2 , error if either are directories
	// 2) len(args) >= 2, *recursive = false, source_file must not be a directory and target must be a directory
	// 3) len(args) >= 2, *recursive = true, target is existing directory
	// 4) len(args) >= 2, *recursive = true, target does not exist

	t, err := os.Stat(target)
	targetExists := (err == nil)
	targetIsDir := (targetExists && t.IsDir())

	if !recursive {
		if len(args) == 2 && !targetIsDir {
			return copyFiles1(sources[0], target, copyfunc)
		} else {
			return copyFiles2(sources, target, copyfunc)
		}
	}

	if targetExists {
		if targetIsDir {
			return copyFiles3(sources, target, copyfunc)
		} else {
			return fmt.Errorf("%s: not a directory", target)
		}
	}

	return copyFiles4(sources, target, copyfunc)
}

func copyFiles1(source string, target string, copyfunc CopyFunc) error {
	if source == target {
		return fmt.Errorf("%s and %s are the same file", source, target)
	}

	s, err := os.Stat(source)

	if err != nil {
		return fmt.Errorf("cannot stat %s: %s", source, err)
	}

	if s.IsDir() {
		return fmt.Errorf("%s is a directory (not copied).", source)
	}

	input, err := os.Open(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer input.Close()

	output, err := os.Create(target)
	if err != nil {
		return err
	}
	defer output.Close()

	return copyfunc(input, output)
}

// expects sources to be files
// expects target to be a directory and to already exist
func copyFiles2(sources []string, target string, copyfunc CopyFunc) error {
	num_err := 0
	for _, source := range sources {
		dest := filepath.Join(target, filepath.Base(source))
		err := copyFiles1(source, dest, copyfunc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			num_err += 1
		}
	}

	if num_err > 0 {
		return fmt.Errorf("non-fatal errors: %d", num_err)
	} else {
		return nil
	}
}

// sources can be files or directories
// expects target to be a directory that already exists
func copyFiles3(sources []string, target string, copyfunc CopyFunc) error {
	num_err := 0
	for _, source := range sources {
		base := filepath.Dir(source)
		err := filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(base, path)
			if err != nil {
				return err
			}

			dest := filepath.Join(target, rel)
			if f.IsDir() {
				os.MkdirAll(dest, os.ModePerm)
				return nil
			}

			return copyFiles1(path, dest, copyfunc)
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			num_err += 1
		}
	}

	if num_err > 0 {
		return fmt.Errorf("non-fatal errors: %d", num_err)
	} else {
		return nil
	}
}

func copyFiles4(sources []string, target string, copyfunc CopyFunc) error {
	fileList := [][2]string{}

	// Assemble list of {sourceFile, sourceDir}
	// Keep the dir so that relative paths may be calculated if needed.
	for _, source := range sources {
		sourceDir := source
		err := filepath.Walk(sourceDir, func(sourceFile string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if f.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(sourceDir, sourceFile)
			if err != nil {
				return err
			}

			targetFile := filepath.Join(target, rel)

			fileList = append(fileList, [2]string{sourceFile, targetFile})
			return nil
		})

		if err != nil {
			return err
		}
	}

	// Copy files
	num_err := 0
	for _, item := range fileList {
		sourceFile := item[0]
		targetFile := item[1]

		os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)
		err := copyFiles1(sourceFile, targetFile, copyfunc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			num_err += 1
		}
	}

	if num_err > 0 {
		return fmt.Errorf("non-fatal errors: %d", num_err)
	} else {
		return nil
	}
}
