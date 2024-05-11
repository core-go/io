package formatter

import (
	"context"
	"github.com/core-go/io/writer"
	"reflect"
)

type FixedLengthFormatter[T any] struct {
	formatCols map[int]*writer.FixedLength
}

func NewFixedLengthFormatter[T any]() (*FixedLengthFormatter[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := writer.GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthFormatter[T]{formatCols: formatCols}, nil
}

func (f *FixedLengthFormatter[T]) Format(ctx context.Context, model *T) string {
	return writer.ToFixedLength(model, f.formatCols)
}
