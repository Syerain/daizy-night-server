package handler

import (
	"github.com/atomreforge/daizy-night-server/internal/service"
)

// HandleXXX() is func of HandlerComplex
type HandlerComplex struct {
	ServiceUser  *service.ServiceUser
	ServiceCode  *service.ServiceCode
	ServiceAdmin *service.ServiceAdmin
}
