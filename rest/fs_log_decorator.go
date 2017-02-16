package rest

import (
	"fmt"
	"net/http"
)

type FileSystemLogDecorator struct {
	FileSystem http.FileSystem
}

func (fs FileSystemLogDecorator) Open(name string) (http.File, error) {
	file, err := fs.FileSystem.Open(name)
	if err != nil {
		fmt.Printf("error while opening file '%v': %v\n", name, err)
	}

	return file, err
}
