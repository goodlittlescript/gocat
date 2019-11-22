package main

import (
	"io"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CopyStream(input io.Reader, output io.Writer, chunk_size int) {
	buffer := make([]byte, chunk_size)
	for {
		nbytes, err := input.Read(buffer)
		if err != nil {
			if err != io.EOF {
				Check(err)
			}

			break
		}

		_, err = output.Write(buffer[:nbytes])
		if err != nil {
			Check(err)
		}
	}
}
