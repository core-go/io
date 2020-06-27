package exporter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

const (
	DateLayout string = "2006-01-02 15:04:05 +0700 +07"
)

type DelimiterFormatter struct {
	Delimiter  string
	modelType  reflect.Type
	formatCols map[int]string
}

func NewDelimiterFormatter(modelType reflect.Type, opts ...string) *DelimiterFormatter {
	sep := ","
	if len(opts) > 0 && len(opts[0]) > 0 {
		sep = opts[0]
	}
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		panic("error get formatCols")
	}
	return &DelimiterFormatter{modelType: modelType, formatCols: formatCols, Delimiter: sep}
}

func (f *DelimiterFormatter) Format(model interface{}) (string, error) {
	arr := make([]string, 0)
	sumValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < sumValue.NumField(); i++ {
		value := fmt.Sprint(sumValue.Field(i).Interface())
		if value == "" || value == "0" || value == "<nil>" {
			value = ""
		} else {
			value = fmt.Sprint(reflect.Indirect(sumValue.Field(i)).Interface())
		}

		if sumValue.Field(i).Type().String() == "string" {
			if strings.Contains(value, f.Delimiter) {
				value = "\"" + value + "\""
			} else {
				if strings.Contains(value, `"`) {
					value = strings.ReplaceAll(value, `"`, `\"`)
				}
			}
		}

		if format, _ := f.formatCols[i]; len(format) > 0 {
			if strings.Contains(format, "dateFormat:") {
				layoutDateStr := strings.ReplaceAll(format, "dateFormat:", "")
				fieldDate, err := time.Parse(DateLayout, value)
				if err != nil {
					fmt.Println("err", fmt.Sprintf("%v", err))
					value = fmt.Sprintf("%v", fmt.Sprintf("%v", value))
				} else {
					value = fmt.Sprintf("%v", fieldDate.UTC().Format(layoutDateStr))
				}
			}
		}
		arr = append(arr, value)
	}
	return strings.Join(arr, f.Delimiter) + "\n", nil
}

func GetIndexesByTag(modelType reflect.Type, tagName string) (map[int]string, error) {
	ma := make(map[int]string, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		if len(tagValue) > 0 {
			ma[i] = tagValue
		} else {
			ma[i] = ""
		}
	}
	return ma, nil
}
