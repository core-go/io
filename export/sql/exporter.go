package sql

import (
	"context"
	"database/sql"
	"reflect"
)

func NewExportRepository(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter, error) {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}
func NewExportAdapter(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter, error) {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}
func NewExportService(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter, error) {
	return NewExporter(db, modelType, buildQuery, transform, write, close)
}

func NewExporter(db *sql.DB, modelType reflect.Type,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, interface{}) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter, error) {
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		return nil, err
	}
	columns := GetColumnsSelect(modelType)
	return &Exporter{DB: db, modelType: modelType, Write: write, Close: close, columns: columns, fieldsIndex: fieldsIndex, Transform: transform, BuildQuery: buildQuery}, nil
}

type Exporter struct {
	DB          *sql.DB
	modelType   reflect.Type
	fieldsIndex map[string]int
	columns     []string
	Transform   func(context.Context, interface{}) string
	BuildQuery  func(context.Context) (string, []interface{})
	Write       func(p []byte) (n int, err error)
	Close       func() error
}

func (s *Exporter) Export(ctx context.Context) (int64, error) {
	query, p := s.BuildQuery(ctx)
	rows, err := s.DB.QueryContext(ctx, query, p...)
	if err != nil {
		return 0, err
	}
	return s.ScanAndWrite(ctx, rows, s.modelType)
}

func (s *Exporter) ScanAndWrite(ctx context.Context, rows *sql.Rows, structType reflect.Type) (int64, error) {
	defer s.Close()

	var i int64
	i = 0
	for rows.Next() {
		initModel := reflect.New(structType).Interface()
		r, swapValues := StructScan(initModel, s.columns, s.fieldsIndex, nil)
		if err := rows.Scan(r...); err != nil {
			return i, err
		}
		SwapValuesToBool(initModel, &swapValues)
		err1 := s.TransformAndWrite(ctx, s.Write, initModel)
		if err1 != nil {
			return i, err1
		}
		i = i + 1
	}
	return i, nil
}

func (s *Exporter) TransformAndWrite(ctx context.Context, write func(p []byte) (n int, err error), model interface{}) error {
	line := s.Transform(ctx, model)
	_, er := write([]byte(line))
	return er
}
