package gocat

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
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
