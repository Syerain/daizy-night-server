package abstract

type InterfaceAppError interface {
	error
	//Say() string
	HttpAbort() int
}
