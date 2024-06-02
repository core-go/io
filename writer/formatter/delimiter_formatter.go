package formatter

import (
	"context"
	"reflect"

	"github.com/core-go/io/writer"
)

type DelimiterFormatter[T any] struct {
	Delimiter  string
	formatCols map[int]writer.Delimiter
}

func NewDelimiterFormatter[T any](opts ...string) (*DelimiterFormatter[T], error) {
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
	return &DelimiterFormatter[T]{formatCols: formatCols, Delimiter: sep}, nil
}

func (f *DelimiterFormatter[T]) Format(ctx context.Context, model *T) string {
	return writer.ToTextWithDelimiter(ctx, model, f.Delimiter, f.formatCols)
}
