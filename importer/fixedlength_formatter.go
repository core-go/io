package importer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewFixedLengthFormatter(modelType reflect.Type) *FixedLengthFormatter {
	formatCols, err := GetIndexesByTag(modelType, "format")
	if err != nil {
		panic("error get formatCols")
	}
	return &FixedLengthFormatter{modelType: modelType, formatCols: formatCols}
}

type FixedLengthFormatter struct {
	modelType  reflect.Type
	formatCols map[int]string
}

func (f FixedLengthFormatter) StringToStruct(ctx context.Context, lines []string) (interface{}, error) {
	line := strings.Join(lines, ``)
	record := reflect.New(f.modelType).Interface()
	err := ScanLineFixLength(line, f.modelType, record, f.formatCols)
	if err != nil {
		return nil, err
	}
	if record != nil {
		return reflect.Indirect(reflect.ValueOf(record)).Interface(), nil
	}
	return record, err
}

func ScanLineFixLength(line string, modelType reflect.Type, record interface{}, formatCols map[int]string) error {
	s := reflect.Indirect(reflect.ValueOf(record))
	numFields := modelType.NumField()
	start := 0
	size := len(line)
	for j := 0; j < numFields; j++ {
		field := modelType.Field(j)
		length, err := strconv.Atoi(field.Tag.Get("length"))
		if err != nil {
			return err
		}
		end := start + length
		if end > size {
			return errors.New(fmt.Sprintf("scanLineFixLength - exceed range max size . Field name = %v , line = %v ", field.Name, line))
		}
		value := line[start:end]
		f := s.Field(j)
		if f.IsValid() {
			if f.CanSet() {
				typef := field.Type.String()
				if f.Kind() == reflect.String {
					stringValue := strings.TrimSpace(value)
					f.SetString(stringValue)
				} else if f.Kind() == reflect.Float64 {
					floatValue, _ := strconv.ParseFloat(value, 64)
					f.SetFloat(floatValue)
				} else if f.Kind() == reflect.Int64 {
					intValue, _ := strconv.ParseInt(value, 64, 0)
					f.SetInt(intValue)
				} else if f.Kind() == reflect.Bool {
					stringValue := strings.TrimSpace(value)
					boolValue, _ := strconv.ParseBool(stringValue)
					f.SetBool(boolValue)
				} else if typef == "*time.Time" || typef == "time.Time" {
					if format, ok := formatCols[j]; ok {
						if strings.Contains(format, "dateFormat:") {
							layoutDateStr := strings.ReplaceAll(format, "dateFormat:", "")
							fieldDate, err := time.Parse(layoutDateStr, value)
							if err != nil {
								return err
							}
							if f.Kind() == reflect.Ptr {
								f.Set(reflect.ValueOf(&fieldDate))
							} else {
								f.Set(reflect.Indirect(reflect.ValueOf(fieldDate)))
							}
						}
					}
				}
			}
		}
		start = end
	}
	return nil
}
