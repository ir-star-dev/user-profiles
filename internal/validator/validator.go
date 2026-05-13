package validator

import (
	"regexp"
	"slices"
	"strings"
)

type Validator struct {
	FieldErrors map[string]string
}

var EmailRX = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@` +
		`[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?` +
		`(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)+$`,
)

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) AddFieldError(name, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[name]; !exists {
		v.FieldErrors[name] = message
	}
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func MaxChars(value string, n int) bool {
	return len(value) <= n
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
