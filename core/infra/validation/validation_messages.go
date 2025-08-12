package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidationMessages(ve validator.ValidationErrors, messages map[string]string) map[string]string {
	msgs := make(map[string]string)
	for _, e := range ve {
		field := CamelToDotNotation(e)
		key := fmt.Sprintf("%s.%s", field, e.Tag())
		if msg, ok := messages[key]; ok {
			msgs[field] = msg
		} else {
			msgs[field] = "Campo obrigatório ou inválido."
		}
	}
	return msgs
}

func CamelToDotNotation(fe validator.FieldError) string {
	ns := fe.StructNamespace()
	parts := strings.Split(ns, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}
	for i, part := range parts {
		parts[i] = toCamelLower(part)
	}
	return strings.Join(parts, ".")
}

func toCamelLower(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
