package v1

import "github.com/atomreforge/daizy-night-server/internal/model"

type RegisterRequest struct {
	model.RegisterBody
}

type RegisterResponse struct {
	Message string `json:"message"`
}
