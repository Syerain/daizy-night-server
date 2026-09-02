package handler

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
)

type InterfaceHandlerComplex interface {
	NewHandlerComplex(
		svcUser abstract.InterfaceServiceUser,
		svcCode abstract.InterfaceServiceCode,
		svcAdmin abstract.InterfaceServiceAdmin,
		svcHealth abstract.InterfaceServiceHealth,
	) *HandlerComplex
}

// HandleXXX() is func of HandlerComplex
type HandlerComplex struct {
	ServiceUser   abstract.InterfaceServiceUser
	ServiceCode   abstract.InterfaceServiceCode
	ServiceAdmin  abstract.InterfaceServiceAdmin
	ServiceHealth abstract.InterfaceServiceHealth
}

func NewHandlerComplex(
	svcUser abstract.InterfaceServiceUser,
	svcCode abstract.InterfaceServiceCode,
	svcAdmin abstract.InterfaceServiceAdmin,
	svcHealth abstract.InterfaceServiceHealth,
) *HandlerComplex {
	return &HandlerComplex{
		ServiceUser:   svcUser,
		ServiceCode:   svcCode,
		ServiceAdmin:  svcAdmin,
		ServiceHealth: svcHealth,
	}
}
