package main

import (
	"html/template"
	"io"
	"os"
)

const (
	_        = iota             // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota) // 1 << (10*1)
	MB                          // 1 << (10*2)
	GB                          // 1 << (10*3)
	TB                          // 1 << (10*4)
	PB                          // 1 << (10*5)
	EB                          // 1 << (10*6)
)

func saveFile(src io.ReadCloser, dest string) error {
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = io.Copy(f, src)
	return err
}

func templateToFile(templateFilename string, filename string, data interface{}) error {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		return err
	}

	return t.Execute(f, data)
}
