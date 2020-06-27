package exporter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FixedLengthFormatter struct {
	sep        string
	modelType  reflect.Type
	formatCols map[int]string
}

func NewFixedLengthFormatWriter(modelType reflect.Type, sep string) *FixedLengthFormatter {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		panic("error get formatCols")
	}
	return &FixedLengthFormatter{modelType: modelType, formatCols: formatCols, sep: sep}
}

func (f *FixedLengthFormatter) Format(model interface{}) (string, error) {
	arr := make([]string, 0)
	sumValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < sumValue.NumField(); i++ {
		value := fmt.Sprint(sumValue.Field(i).Interface())
		field := f.modelType.Field(i)
		length, err := strconv.Atoi(field.Tag.Get("length"))
		if err != nil {
			return "", err
		}
		if value == "" || value == "0" || value == "<nil>" {
			value = ""
		} else if format, _ := f.formatCols[i]; len(format) > 0 && strings.Contains(format, "dateFormat:") {
			layoutDateStr := strings.ReplaceAll(format, "dateFormat:", "")
			fieldDate, err := time.Parse(DateLayout, value)
			if err != nil {
				fmt.Println("err", fmt.Sprintf("%v", err))
				value = fmt.Sprintf("%v", fmt.Sprintf("%v", value))
			} else {
				value = fmt.Sprintf("%v", fieldDate.UTC().Format(layoutDateStr))
			}
		} else {
			if len(value) > length {
				value = strings.TrimSpace(value)
			}
			if format, _ := f.formatCols[i]; len(format) > 0 {
				value = fmt.Sprintf(format, value)
			}
			value = FixedLengthString(length, value)
		}
		arr = append(arr, value)
	}
	return strings.Join(arr, f.sep) + "\n", nil
}

func FixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}
