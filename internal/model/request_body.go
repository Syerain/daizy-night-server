package model

import (
	"time"

	"github.com/atomreforge/daizy-night-server/internal/consts"
)

type RegisterBody struct {
	Registerway  consts.Registerway `json:"registerway"`
	Username     string             `json:"username"`
	Nickname     string             `json:"nickname"`
	Password     string             `json:"password"`
	Registercode RegistercodeRawHex `json:"registercode"`
}

type LoginBody struct {
	Loginway  consts.Loginway `json:"loginway"`
	Username  string          `json:"username"`
	Password  string          `json:"password"`
	Entrycode string          `json:"entrycode"`
}

type InfoMe struct {
	Uid uint `json:"uid"`

	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`

	RegisterTime time.Time `json:"register_time"`

	Role consts.Role `json:"role"`

	GitHubID    *int64 `json:"github_id"`
	GitHubLogin string `json:"github_login"`
}
