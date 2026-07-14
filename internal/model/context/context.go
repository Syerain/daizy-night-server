package model

import "github.com/labstack/echo/v5"

type ContextRegister struct {
	echo.Context
	Username     string
	Nickname     string
	Password     string
	Registercode string
}
