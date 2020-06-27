package export

import (
	"bytes"
	"context"
	"database/sql"
	"reflect"
)
func NewExportRepository(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) (string, error),
	write func(p []byte) (n int, err error),
	close func() error,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}
func NewExportAdapter(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) (string, error),
	write func(p []byte) (n int, err error),
	close func() error,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}
func NewExportService(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) (string, error),
	write func(p []byte) (n int, err error),
	close func() error,
) *Exporter {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}

func NewExporter(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) (string, error),
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
	Transform   func(context.Context, interface{}) (string, error)
	BuildQuery  func(context.Context) (string, []interface{})
	Write       func(p []byte) (n int, err error)
	Close       func() error
}

func (s *Exporter) Export(ctx context.Context) error {
	query, p := s.BuildQuery(ctx)
	rows, err := s.DB.QueryContext(ctx, query, p...)
	if err != nil {
		return err
	}
	return s.ScanAndWrite(ctx, rows, s.modelType)
}

func (s *Exporter) ScanAndWrite(ctx context.Context, rows *sql.Rows, structType reflect.Type) error {
	defer s.Close()

	for rows.Next() {
		initModel := reflect.New(structType).Interface()
		r, swapValues := StructScan(initModel, nil, s.fieldsIndex, nil)
		if err := rows.Scan(r...); err != nil {
			return err
		}
		SwapValuesToBool(initModel, &swapValues)
		err1 := s.TransformAndWrite(ctx, s.Write, initModel)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func (s *Exporter) TransformAndWrite(ctx context.Context, write func(p []byte) (n int, err error), model interface{}) error {
	var buffer bytes.Buffer
	line, err := s.Transform(ctx, model)
	if err != nil {
		return err
	}
	buffer.WriteString(line)

	_, err0 := write(buffer.Bytes())
	return err0
}
