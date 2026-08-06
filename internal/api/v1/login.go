package v1

import (
	"github.com/atomreforge/daizy-night-server/internal/model"
)

type LoginRequest struct {
	model.LoginBody
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
