package service

import "log/slog"

type ServiceAdmin struct {
}

func NewServiceAdmin() *ServiceAdmin {
	return &ServiceAdmin{}
}

func (*ServiceAdmin) Sudo() {
	slog.Info("admin executed sudo")
}
