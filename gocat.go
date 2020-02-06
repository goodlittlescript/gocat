package gocat

import (
	"io"
)

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
