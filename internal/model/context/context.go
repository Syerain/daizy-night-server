package model

import "github.com/labstack/echo/v5"

// not used yet, probably never ?
type ContextRegister struct {
	echo.Context
	Username     string
	Nickname     string
	Password     string
	Registercode string
}
