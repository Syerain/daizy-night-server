package consts

type Registerway string
type Loginway string

const (
	RegisterLegacy Registerway = "legacy"
	RegisterGithub Registerway = "oauth-github"
)

const (
	LoginLegacy Loginway = "legacy"
	LoginGithub Loginway = "oauth-github"
)
