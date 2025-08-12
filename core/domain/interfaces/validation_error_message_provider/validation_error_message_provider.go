package validationerrormessageprovider

import "github.com/go-playground/validator/v10"

type ValidationMessageProvider interface {
	ValidationMessages(validator.ValidationErrors) map[string]string
}
