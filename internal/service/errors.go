package service

// login
type ErrUserLogin struct {
	User    string
	Message string
}

func (e *ErrUserLogin) Error() string {
	return "User login error: " + e.User + " - " + e.Message
}

// regiester
type ErrRegister struct {
	Message string
}

func (e *ErrRegister) Error() string {
	return "Register error: " + e.Message
}
