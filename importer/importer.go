package importer

import (
	"database/sql"
	"io"
	"reflect"
)

func NewImporter(db *sql.DB, modelType reflect.Type,
	transform func(lines []string) (interface{}, error),
	write func(data interface{}, endLineFlag bool) error,
	read func(next func(lines []string, err error) error) error,
) *Importer {
	return &Importer{DB: db, modelType: modelType, Transform: transform, Write: write, Read: read}
}

type Importer struct {
	DB        *sql.DB
	modelType reflect.Type
	Transform func(lines []string) (interface{}, error)
	Read      func(next func(lines []string, err error) error) error
	Write     func(data interface{}, endLineFlag bool) error
}

func (s *Importer) Import() (err error) {
	err = s.Read(func(lines []string, err error) error {
		if err == io.EOF {
			err = s.Write(nil, true)
			return nil
		}
		itemStruct, err := s.Transform(lines)
		if err != nil {
			return err
		}
		err = s.Write(itemStruct, false)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}
