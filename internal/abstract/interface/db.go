package abstract

type InterfaceStatus interface {
	Check() error
	Close() error
}
