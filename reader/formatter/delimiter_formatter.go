package formatter

import (
	"context"
	"reflect"
	"strings"

	"github.com/core-go/io/reader"
)

func NewDelimiterFormatter[T any](options ...string) (*DelimiterFormatter[T], error) {
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
	return &DelimiterFormatter[T]{formatCols: formatCols, separator: separator}, nil
}

type DelimiterFormatter[T any] struct {
	formatCols map[int]reader.Delimiter
	separator  string
}

func (f DelimiterFormatter[T]) ToStruct(ctx context.Context, lineStr string) (T, error) {
	lines := strings.Split(lineStr, f.separator)
	var res T
	err := reader.ScanLine(lines, &res, f.formatCols)
	return res, err
}
