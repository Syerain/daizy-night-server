package confx

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validatable interface {
	Validate() error
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	v.RegisterTagNameFunc(resolveKeyName)
	return &CustomValidator{v}
}

type CustomValidator struct {
	validate *validator.Validate
}

func (v *CustomValidator) Validate(i any) error {
	if sv, ok := i.(Validatable); ok {
		return sv.Validate()
	}

	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return fmt.Errorf("invalid validation error: %w", err)
	}

	validationErrors := err.(validator.ValidationErrors)
	translatedErrors := make(map[string]string)

	for _, e := range validationErrors {
		// 处理 Namespace (去除结构体前缀)
		namespace := e.Namespace()
		if i := strings.Index(namespace, "."); i != -1 {
			namespace = namespace[i+1:]
		}

		if e.Param() != "" {
			translatedErrors[namespace] = fmt.Sprintf("%s=%s", e.Tag(), e.Param())
		} else {
			translatedErrors[namespace] = e.Tag()
		}
	}

	return &ValidationError{Errors: translatedErrors}
}
