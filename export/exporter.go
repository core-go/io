package export

import (
	"context"
	"database/sql"
	"reflect"
)

func NewExportRepository[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter[T], error) {
	return NewExporter[T](db, buildQuery, transform, write, close)
}
func NewExportAdapter[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter[T], error) {
	return NewExporter[T](db, buildQuery, transform, write, close)
}
func NewExportService[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter[T], error) {
	return NewExporter[T](db, buildQuery, transform, write, close)
}

func NewExporter[T any](db *sql.DB,
	buildQuery func(context.Context) (string, []interface{}),
	transform func(context.Context, *T) string,
	write func(p []byte) (n int, err error),
	close func() error,
) (*Exporter[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	fieldsIndex, err := GetColumnIndexes(modelType)
	if err != nil {
		return nil, err
	}
	columns := GetColumnsSelect(modelType)
	return &Exporter[T]{DB: db, Write: write, Close: close, columns: columns, fieldsIndex: fieldsIndex, Transform: transform, BuildQuery: buildQuery}, nil
}

type Exporter[T any] struct {
	DB          *sql.DB
	fieldsIndex map[string]int
	columns     []string
	Transform   func(context.Context, *T) string
	BuildQuery  func(context.Context) (string, []interface{})
	Write       func(p []byte) (n int, err error)
	Close       func() error
}

func (s *Exporter[T]) Export(ctx context.Context) (int64, error) {
	query, p := s.BuildQuery(ctx)
	rows, err := s.DB.QueryContext(ctx, query, p...)
	if err != nil {
		return 0, err
	}
	return s.ScanAndWrite(ctx, rows)
}

func (s *Exporter[T]) ScanAndWrite(ctx context.Context, rows *sql.Rows) (int64, error) {
	defer s.Close()

	var i int64
	i = 0
	for rows.Next() {
		var obj T
		r, swapValues := StructScan(&obj, s.columns, s.fieldsIndex, nil)
		if err := rows.Scan(r...); err != nil {
			return i, err
		}
		SwapValuesToBool(&obj, &swapValues)
		err1 := s.TransformAndWrite(ctx, s.Write, &obj)
		if err1 != nil {
			return i, err1
		}
		i = i + 1
	}
	return i, nil
}

func (s *Exporter[T]) TransformAndWrite(ctx context.Context, write func(p []byte) (n int, err error), model *T) error {
	line := s.Transform(ctx, model)
	_, er := write([]byte(line))
	return er
}
