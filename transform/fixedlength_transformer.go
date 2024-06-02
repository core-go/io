package transform

import (
	"context"
	"reflect"

	"github.com/core-go/io/reader"
)

func NewFixedLengthTransformer[T any]() (*FixedLengthTransformer[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := reader.GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthTransformer[T]{formatCols: formatCols}, nil
}

type FixedLengthTransformer[T any] struct {
	formatCols map[int]*reader.FixedLength
}

func (f FixedLengthTransformer[T]) Transform(ctx context.Context, line string) (T, error) {
	var res T
	err := reader.ScanLineFixLength(line, &res, f.formatCols)
	return res, err
}
