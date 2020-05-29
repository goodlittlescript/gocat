package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	filepath.Walk("test", func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Println(path)

		return nil
	})
}
