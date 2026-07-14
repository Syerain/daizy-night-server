package model

import (
	"daizynight/internal/constants"
)

type RegisterBody struct {
	Username     string
	Nickname     string
	Password     string
	Registercode string
	Registerway  constants.RegisterWay
}
