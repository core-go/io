package formatter

import (
	"context"
	"reflect"

	"github.com/core-go/io/writer"
)

type DelimiterTransformer[T any] struct {
	Delimiter  string
	formatCols map[int]writer.Delimiter
}

func NewDelimiterTransformer[T any](opts ...string) (*DelimiterTransformer[T], error) {
	sep := "|"
	if len(opts) > 0 && len(opts[0]) > 0 {
		sep = opts[0]
	}
	skipTag := ""
	if len(opts) > 1 && len(opts[1]) > 0 {
		skipTag = opts[1]
	}
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := writer.GetIndexesByTag(modelType, "format", skipTag)
	if err != nil {
		return nil, err
	}
	return &DelimiterTransformer[T]{formatCols: formatCols, Delimiter: sep}, nil
}
func NewDelimiterFormatter[T any](opts ...string) (*DelimiterTransformer[T], error) {
	return NewDelimiterTransformer[T](opts...)
}

func (f *DelimiterTransformer[T]) Transform(ctx context.Context, model *T) string {
	return writer.ToTextWithDelimiter(model, f.Delimiter, f.formatCols)
}
