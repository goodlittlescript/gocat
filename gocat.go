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

type CatFunc func(io.Reader, io.Writer) error
type CopyFunc func(string, string) error
type FailFunc func(error)

func NewCopyFunc(catfunc CatFunc) CopyFunc {
	return func(source string, target string) error {
		input, err := os.Open(source)
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := os.Create(target)
		if err != nil {
			return err
		}
		defer output.Close()

		return catfunc(input, output)
	}
}

func CatFiles(files []string, output io.Writer, catfunc CatFunc) error {
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

		err = catfunc(input, output)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyFiles(args []string, recursive bool, copyfunc CopyFunc, failfunc FailFunc) {
	err := copyFiles(args, recursive, copyfunc, failfunc)
	if err != nil {
		failfunc(err)
	}
}

func copyFiles(args []string, recursive bool, copyfunc CopyFunc, failfunc FailFunc) error {
	if len(args) == 0 {
		return fmt.Errorf("missing file operand")
	}

	if len(args) == 1 {
		return fmt.Errorf("missing destination file operand after '%s'", args[0])
	}

	sources := args[:len(args)-1]
	target := args[len(args)-1]

	t, err := os.Stat(target)
	targetExists := (err == nil)
	targetIsDir := (targetExists && t.IsDir())

	if !recursive {
		if len(args) == 2 && !targetIsDir {
			return copyFiles1(sources[0], target, copyfunc)
		} else {
			return copyFiles2(sources, target, copyfunc, failfunc)
		}
	}

	if targetExists {
		if !targetIsDir {
			return fmt.Errorf("%s: not a directory", target)
		}
	} else {
		if len(sources) > 1 {
			return fmt.Errorf("target '%s' is not a directory", target)
		}
	}

	return copyFiles3(sources, target, targetExists, copyfunc, failfunc)
}

// First Synopsis Form:
//
// The first synopsis form is denoted by two operands, neither of which are
// existing files of type directory. The cp utility shall copy the contents of
// source_file (or, if source_file is a file of type symbolic link, the contents
// of the file referenced by source_file) to the destination path named by
// target_file.
//
// Note: care is taken in the calling contexts to check that target is not a
// directory, but herein source must be checked to not be a directory.
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

	return copyfunc(source, target)
}

// Second Synopsis Form:
//
// The second synopsis form is denoted by two or more operands where the -R or
// -r options are not specified and the first synopsis form is not applicable.
// It shall be an error if any source_file is a file of type directory, if
// target does not exist, or if target is a file of a type defined by the System
// Interfaces volume of IEEE Std 1003.1-2001, but is not a file of type
// directory. The cp utility shall copy the contents of each source_file (or, if
// source_file is a file of type symbolic link, the contents of the file
// referenced by source_file) to the destination path named by the concatenation
// of target, a slash character, and the last component of source_file.
func copyFiles2(sources []string, target string, copyfunc CopyFunc, failfunc FailFunc) error {
	for _, source := range sources {
		dest := filepath.Join(target, filepath.Base(source))
		err := copyFiles1(source, dest, copyfunc)
		if err != nil {
			failfunc(err)
		}
	}

	return nil
}

// Third/Fourth Synopsis Form:
//
// The third and fourth synopsis forms are denoted by two or more operands where
// the -R or -r options are specified. The cp utility shall copy each file in
// the file hierarchy rooted in each source_file to a destination path named as
// follows:
//
// * If target exists and is a file of type directory, the name of the
// corresponding destination path for each file in the file hierarchy shall be
// the concatenation of target, a slash character, and the pathname of the file
// relative to the directory containing source_file.
//
// * If target does not exist and two operands are specified, the name of the
// corresponding destination path for source_file shall be target; the name of
// the corresponding destination path for all other files in the file hierarchy
// shall be the concatenation of target, a slash character, and the pathname of
// the file relative to source_file.
//
// It shall be an error if target does not exist and more than two operands are
// specified, or if target exists and is a file of a type defined by the System
// Interfaces volume of IEEE Std 1003.1-2001, but is not a file of type
// directory.
func copyFiles3(sources []string, target string, targetExists bool, copyfunc CopyFunc, failfunc FailFunc) error {
	for _, source := range sources {
		filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			rel, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			var dest string
			if targetExists {
				dest = filepath.Join(target, filepath.Base(source), rel)
			} else {
				dest = filepath.Join(target, rel)
			}

			if f.IsDir() {
				os.MkdirAll(dest, os.ModePerm)
				return nil
			}

			err = copyFiles1(path, dest, copyfunc)
			if err != nil {
				failfunc(err)
			}
			return nil
		})
	}

	return nil
}
