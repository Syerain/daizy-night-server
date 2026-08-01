package model

import (
	"github.com/atomreforge/dnserver/internal/constants"
)

type RegisterBody struct {
	Username     string
	Nickname     string
	Password     string
	Registercode string
	Registerway  constants.Registerway
}

type LoginBody struct {
	Loginway  constants.Loginway
	Username  string
	Atomid    int
	Password  string
	Entrycode string
}
