package reader

import (
	"context"
	"errors"
	"golang.org/x/text/encoding"
	"reflect"
)

type FileType string

const (
	DelimiterType   FileType = "Delimiter"
	FixedlengthType FileType = "Fixedlength"
)

type Formater interface {
	ToStruct(ctx context.Context, line string, res interface{}) error
}

func NewFormater(modelType reflect.Type, fileType FileType, opts ...string) (Formater, error) {
	if fileType == DelimiterType {
		return NewDelimiterFormatter(modelType, opts...)
	}
	if fileType == FixedlengthType {
		return NewFixedLengthFormatter(modelType)
	}
	return nil, errors.New("Bad csv type")
}

type Reader interface {
	Read(next func(lines string, err error, numLine int) error) error
}

func NewReader(buildFileName func() string, csvType FileType, opts ...*encoding.Decoder) (Reader, error) {
	if csvType == DelimiterType {
		return NewDelimiterFileReader(buildFileName)
	}
	if csvType == FixedlengthType {
		return NewFixedlengthFileReader(buildFileName, opts...)
	}
	return nil, errors.New("Bad csv type")
}
