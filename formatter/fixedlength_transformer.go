package formatter

import (
	"context"
	"github.com/core-go/io/writer"
	"reflect"
)

type FixedLengthTransformer[T any] struct {
	formatCols map[int]*writer.FixedLength
}

func NewFixedLengthTransformer[T any]() (*FixedLengthTransformer[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := writer.GetIndexes(modelType, "format")
	if err != nil {
		return nil, err
	}
	return &FixedLengthTransformer[T]{formatCols: formatCols}, nil
}

func NewFixedLengthFormatter[T any]() (*FixedLengthTransformer[T], error) {
	return NewFixedLengthTransformer[T]()
}
func (f *FixedLengthTransformer[T]) Transform(ctx context.Context, model *T) string {
	return writer.ToFixedLength(model, f.formatCols)
}
