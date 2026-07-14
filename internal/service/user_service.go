package service

import (
	"daizynight/internal/model"
)

func Register(b *model.RegisterBody) error {
	err := ValidateRegisterParams(b)
	if err != nil {
		return err
	}
	return nil
}
