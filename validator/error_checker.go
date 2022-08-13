package validator

import sv "github.com/core-go/io/import"

func NewErrorChecker() *sv.ErrorChecker {
	v := NewValidator()
	return sv.NewErrorChecker(v.Validate)
}
