package service

// login
type UserLoginError struct {
	User    string
	Message string
}

func (e *UserLoginError) Error() string {
	return "User login error: " + e.User + " - " + e.Message
}

// regiester
type RegisterError struct {
	Message string
}

func (e *RegisterError) Error() string {
	return "Regiester error: " + e.Message
}
