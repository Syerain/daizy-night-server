package model

import (
	"daizynight/internal/constants"
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
	Atomid    string
	Password  string
	Entrycode string
}
