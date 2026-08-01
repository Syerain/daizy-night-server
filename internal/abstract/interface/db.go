package abstract

type InterfaceStatus interface {
	Check() error
	Close() error

	// already stated by repo interface
	/* CreateUser(b *model.User) error
	GetUserByUsername(name string) (*model.User, error)
	GetUserByAtomid(atomid uint) (*model.User, error)
	SaveRefreshToken(atomid int, rawToken string) error
	GetRefreshToken(atomid int, rawToken string) (bool, error)
	RevokeUserTokens(atomid int) error */
}
