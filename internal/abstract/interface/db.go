package abstract

import "gorm.io/gorm"

type InterfaceStatus interface {
	Check() error
	Close() error
	Reboot() error
}

type InterfaceProviderDB interface {
	Check() error
	Close() error
	DB() *gorm.DB
}
