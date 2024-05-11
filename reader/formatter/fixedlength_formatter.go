package formatter

import (
	"context"
	"reflect"

	"github.com/core-go/io/reader"
)

func NewFixedLengthFormatter[T any]() (*FixedLengthFormatter[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := reader.GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthFormatter[T]{formatCols: formatCols}, nil
}

type FixedLengthFormatter[T any] struct {
	formatCols map[int]*reader.FixedLength
}

func (f FixedLengthFormatter[T]) ToStruct(ctx context.Context, line string) (T, error) {
	var res T
	err := reader.ScanLineFixLength(line, &res, f.formatCols)
	return res, err
}
