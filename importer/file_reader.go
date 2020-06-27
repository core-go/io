package importer

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strings"
)

type FileReader struct {
	FileName string
}

func NewFileReader(buildFileName func() string) (*FileReader, error) {
	var fr FileReader
	fileName := buildFileName()
	if len(strings.TrimSpace(fileName)) == 0 {
		return nil, errors.New("file name cannot be empty")
	}
	fr.FileName = fileName
	return &fr, nil
}

func (fr *FileReader) Read(next func(lines []string, err error) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		err = errors.New("cannot open file")
		next(make([]string, 0), err)
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() != "" {
			err := next([]string{scanner.Text()}, nil)
			if err != nil {
				return err
			}
		}
	}
	next([]string{}, io.EOF)
	return nil
}

func (fr *FileReader) ReadFileCSV(next func(lines []string, err error) error) error {
	file, err := os.Open(fr.FileName)
	if err != nil {
		next(make([]string, 0), err)
	}
	// Create a new reader.
	r := csv.NewReader(file)
	defer file.Close()
	for {
		record, err := r.Read()
		err2 := next(record, err)
		if err2 != nil {
			return err2
		}
		if err == io.EOF {
			break
		}
	}
	return err
}
