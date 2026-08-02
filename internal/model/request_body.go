package model

import (
	"github.com/atomreforge/daizy-night-server/internal/consts"
)

type RegisterBody struct {
	Registerway  consts.Registerway
	Username     string
	Nickname     string
	Password     string
	Registercode string
}

type LoginBody struct {
	Loginway  consts.Loginway
	Atomid    uint
	Username  string
	Password  string
	Entrycode string
}
