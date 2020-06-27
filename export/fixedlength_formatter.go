package export

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FixedLengthFormatter struct {
	modelType  reflect.Type
	formatCols map[int]*FixedLength
}
type FixedLength struct {
	Format string
	Length int
}
func GetIndexes(modelType reflect.Type, tagName string) (map[int]*FixedLength, error) {
	ma := make(map[int]*FixedLength, 0)
	if modelType.Kind() != reflect.Struct {
		return ma, errors.New("bad type")
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tagValue := field.Tag.Get(tagName)
		tagLength := field.Tag.Get("length")
		if len(tagLength) > 0 {
			length, err := strconv.Atoi(tagLength)
			if err != nil || length < 0 {
				return ma, err
			}
			v := &FixedLength{Length: length}
			if len(tagValue) > 0 {
				v.Format = tagValue
			}
			ma[i] = v
		}
	}
	return ma, nil
}

func NewFixedLengthFormatter(modelType reflect.Type) *FixedLengthFormatter {
	formatCols, err := GetIndexes(modelType, "format")
	if err != nil {
		panic("error get formatCols")
	}
	return &FixedLengthFormatter{modelType: modelType, formatCols: formatCols}
}

func (f *FixedLengthFormatter) Format(ctx context.Context, model interface{}) (string, error) {
	return ToFixedLength(model, f.formatCols)
}
func ToFixedLength(model interface{}, formatCols map[int]*FixedLength) (string, error) {
	arr := make([]string, 0)
	sumValue := reflect.Indirect(reflect.ValueOf(model))
	for i := 0; i < sumValue.NumField(); i++ {
		format, ok := formatCols[i]
		if ok {
			value := fmt.Sprint(sumValue.Field(i).Interface())
			if value == "" || value == "0" || value == "<nil>" {
				value = ""
			} else if len(format.Format) > 0 && strings.Contains(format.Format, "dateFormat:") {
				layoutDateStr := strings.ReplaceAll(format.Format, "dateFormat:", "")
				fieldDate, err := time.Parse(DateLayout, value)
				if err != nil {
					fmt.Println("err", fmt.Sprintf("%v", err))
					value = fmt.Sprintf("%v", fmt.Sprintf("%v", value))
				} else {
					value = fmt.Sprintf("%v", fieldDate.UTC().Format(layoutDateStr))
				}
			} else {
				if len(value) > format.Length {
					value = strings.TrimSpace(value)
				}
				if len(format.Format) > 0 {
					value = fmt.Sprintf(format.Format, value)
				}
				value = FixedLengthString(format.Length, value)
			}
			arr = append(arr, value)
		}
	}
	return strings.Join(arr, "") + "\n", nil
}
func FixedLengthString(length int, str string) string {
	verb := fmt.Sprintf("%%%d.%ds", length, length)
	return fmt.Sprintf(verb, str)
}
