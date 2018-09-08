package main

import (
	"io"
	"os"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func CopyStream(input *os.File, output *os.File, chunk_size int) {
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
