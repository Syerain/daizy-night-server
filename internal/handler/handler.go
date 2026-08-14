package handler

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
)

type InterfaceHandlerComplex interface {
	NewHandlerComplex(
		svcUser abstract.InterfaceServiceUser,
		svcCode abstract.InterfaceServiceCode,
		svcAdmin abstract.InterfaceServiceAdmin,
	) *HandlerComplex
}

// HandleXXX() is func of HandlerComplex
type HandlerComplex struct {
	ServiceUser  abstract.InterfaceServiceUser
	ServiceCode  abstract.InterfaceServiceCode
	ServiceAdmin abstract.InterfaceServiceAdmin
}

func NewHandlerComplex(svcUser abstract.InterfaceServiceUser, svcCode abstract.InterfaceServiceCode, svcAdmin abstract.InterfaceServiceAdmin) *HandlerComplex {
	return &HandlerComplex{
		ServiceUser:  svcUser,
		ServiceCode:  svcCode,
		ServiceAdmin: svcAdmin,
	}
}
