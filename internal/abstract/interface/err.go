package abstract

type InterfaceAppError interface {
	error
	//Say() string
	StatusCode() int
	Respond() string
}
