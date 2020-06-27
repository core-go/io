package exporter

import (
	"bytes"
	"database/sql"
	"reflect"
)

func NewExporter(db *sql.DB, modelType reflect.Type,
	buildQuery func() (string, []interface{}),
	transform func(model interface{}) (string, error),
	write func(p []byte) (n int, err error),
	close func() error,
) *Exporter {
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		panic("error get fieldsIndex")
	}
	columns := GetColumnsSelect(modelType)
	return &Exporter{DB: db, modelType: modelType, Write: write, Close: close, columns: columns, fieldsIndex: fieldsIndex, Transform: transform, BuildQuery: buildQuery}
}

type Exporter struct {
	DB          *sql.DB
	modelType   reflect.Type
	fieldsIndex map[string]int
	columns     []string
	Transform   func(model interface{}) (string, error)
	BuildQuery  func() (string, []interface{})
	Write       func(p []byte) (n int, err error)
	Close       func() error
}

func (s *Exporter) Export() error {
	query, p := s.BuildQuery()
	rows, err := s.DB.Query(query, p...)
	if err != nil {
		return err
	}
	return s.ScanAndWrite(rows, s.modelType)
}

func (s *Exporter) ScanAndWrite(rows *sql.Rows, structType reflect.Type) error {
	defer s.Close()

	for rows.Next() {
		initModel := reflect.New(structType).Interface()
		r, swapValues := StructScan(initModel, nil, s.fieldsIndex, nil)
		if err := rows.Scan(r...); err != nil {
			return err
		}
		SwapValuesToBool(initModel, &swapValues)
		err1 := s.TransformAndWrite(s.Write, initModel)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func (s *Exporter) TransformAndWrite(write func(p []byte) (n int, err error), model interface{}) error {
	var buffer bytes.Buffer
	line, err := s.Transform(model)
	if err != nil {
		return err
	}
	buffer.WriteString(line)

	_, err0 := write(buffer.Bytes())
	return err0
}
