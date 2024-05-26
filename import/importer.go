package importer

import (
	"context"
	"fmt"
	"io"
	"reflect"
)

type ErrorMessage struct {
	Field   string `yaml:"field" mapstructure:"field" json:"field,omitempty" gorm:"column:field" bson:"field,omitempty" dynamodbav:"field,omitempty" firestore:"field,omitempty"`
	Code    string `yaml:"code" mapstructure:"code" json:"code,omitempty" gorm:"column:code" bson:"code,omitempty" dynamodbav:"code,omitempty" firestore:"code,omitempty"`
	Param   string `yaml:"param" mapstructure:"param" json:"param,omitempty" gorm:"column:param" bson:"param,omitempty" dynamodbav:"param,omitempty" firestore:"param,omitempty"`
	Message string `yaml:"message" mapstructure:"message" json:"message,omitempty" gorm:"column:message" bson:"message,omitempty" dynamodbav:"message,omitempty" firestore:"message,omitempty"`
}

type ErrorHandler struct {
	HandleError func(ctx context.Context, format string, fields map[string]interface{})
	FileName    string
	LineNumber  string
	Map         *map[string]interface{}
}

func NewErrorHandler(logger func(ctx context.Context, format string, fields map[string]interface{}), fileName string, lineNumber string, mp *map[string]interface{}) *ErrorHandler {
	if len(fileName) <= 0 {
		fileName = "filename"
	}
	if len(lineNumber) <= 0 {
		lineNumber = "lineNumber"
	}
	return &ErrorHandler{
		HandleError: logger,
		FileName:    fileName,
		LineNumber:  lineNumber,
		Map:         mp,
	}
}

func (e *ErrorHandler) HandlerError(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v", raw, rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v line: %d", raw, rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v filename:%s", raw, rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Message is invalid: %s %+v . Error: %+v filename:%s line: %d", raw, rs, err, fileName, i), ext)
	}
}

func (e *ErrorHandler) HandlerException(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string) {
	var ext = make(map[string]interface{})
	if e.Map != nil {
		ext = *e.Map
	}
	if len(e.FileName) > 0 && len(e.LineNumber) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v", raw, rs, err), ext)
	} else if len(e.FileName) > 0 {
		if len(fileName) > 0 {
			ext[e.FileName] = fileName
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v line: %d", raw, rs, err, i), ext)
	} else if len(e.LineNumber) > 0 {
		if i > 0 {
			ext[e.LineNumber] = i
		}
		e.HandleError(ctx, fmt.Sprintf("Error to write: %s %+v . Error: %+v filename:%s", raw, rs, err, fileName), ext)
	} else {
		e.HandleError(ctx, fmt.Sprintf("Error to write:  %s %+v . Error: %v filename: %s line: %d", raw, rs, err, fileName, i), ext)
	}
}

func NewImportRepository(modelType reflect.Type,
	transform func(ctx context.Context, lines string, res interface{}) error,
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data interface{}) error,
	opt ...func(ctx context.Context) error,
) *Importer {
	return NewImporter(modelType, transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImportAdapter(modelType reflect.Type,
	transform func(ctx context.Context, lines string, res interface{}) error,
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data interface{}) error,
	opt ...func(ctx context.Context) error,
) *Importer {
	return NewImporter(modelType, transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImportService(modelType reflect.Type,
	transform func(ctx context.Context, lines string, res interface{}) error,
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data interface{}) error,
	opt ...func(ctx context.Context) error,
) *Importer {
	return NewImporter(modelType, transform, read, handleException, validate, handleError, filename, write, opt...)
}
func NewImporter(modelType reflect.Type,
	transform func(ctx context.Context, lines string, res interface{}) error,
	read func(next func(line string, err error, numLine int) error) error,
	handleException func(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string),
	validate func(ctx context.Context, model interface{}) ([]ErrorMessage, error),
	handleError func(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string),
	filename string,
	write func(ctx context.Context, data interface{}) error,
	opt ...func(ctx context.Context) error,
) *Importer {
	var flush func(ctx context.Context) error
	if len(opt) > 0 {
		flush = opt[0]
	}
	return &Importer{modelType: modelType, Transform: transform, Write: write, Flush: flush, Read: read, Validate: validate, HandleError: handleError, HandleException: handleException, Filename: filename}
}

type Importer struct {
	modelType       reflect.Type
	Transform       func(ctx context.Context, lines string, res interface{}) error
	Read            func(next func(line string, err error, numLine int) error) error
	Validate        func(ctx context.Context, model interface{}) ([]ErrorMessage, error)
	HandleError     func(ctx context.Context, raw string, rs interface{}, err []ErrorMessage, i int, fileName string)
	HandleException func(ctx context.Context, raw string, rs interface{}, err error, i int, fileName string)
	Filename        string
	Write           func(ctx context.Context, data interface{}) error
	Flush           func(ctx context.Context) error
}

func (s *Importer) Import(ctx context.Context) (total int, success int, err error) {
	err = s.Read(func(line string, err error, numLine int) error {
		if err == io.EOF {
			if s.Flush != nil {
				return s.Flush(ctx)
			}
			return nil
		}
		total++
		record := reflect.New(s.modelType).Interface()
		err = s.Transform(ctx, line, record)
		if err != nil {
			return err
		}
		if s.Validate != nil {
			errs, err := s.Validate(ctx, record)
			if err != nil {
				return err
			}
			if len(errs) > 0 {
				s.HandleError(ctx, line, record, errs, numLine, s.Filename)
				return nil
			}
		}
		err = s.Write(ctx, record)
		if err != nil {
			if s.HandleException != nil {
				s.HandleException(ctx, line, record, err, numLine, s.Filename)
				return nil
			} else {
				return err
			}
		}
		success++
		return nil
	})
	if err != nil && err != io.EOF {
		return total, success, err
	}
	return total, success, nil
}
