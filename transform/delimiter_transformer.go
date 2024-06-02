package transform

import (
	"context"
	"reflect"
	"strings"

	"github.com/core-go/io/reader"
)

func NewDelimiterTransformer[T any](options ...string) (*DelimiterTransformer[T], error) {
	var t T
	modelType := reflect.TypeOf(t)
	formatCols, err := reader.GetIndexesByTag(modelType, "format")
	if err != nil {
		return nil, err
	}
	separator := ""
	if len(options) > 0 {
		separator = options[0]
	} else {
		separator = "|"
	}
	return &DelimiterTransformer[T]{formatCols: formatCols, separator: separator}, nil
}

type DelimiterTransformer[T any] struct {
	formatCols map[int]reader.Delimiter
	separator  string
}

func (f DelimiterTransformer[T]) ToStruct(ctx context.Context, lineStr string) (T, error) {
	lines := strings.Split(lineStr, f.separator)
	var res T
	err := reader.ScanLine(lines, &res, f.formatCols)
	return res, err
}
