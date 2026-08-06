package service

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

var _ abstract.InterfaceServiceCode = (*ServiceCode)(nil)

type ServiceCode struct {
	RepoRegistercode abstract.InterfaceRepoRegistercode
}

func NewServiceCode(repoRegcode abstract.InterfaceRepoRegistercode) *ServiceCode {
	return &ServiceCode{
		RepoRegistercode: repoRegcode,
	}
}

func (s *ServiceCode) RecordRegistercode(rawHex model.RegistercodeRawHex) error {
	if err := s.RepoRegistercode.Record(rawHex); err != nil {
		return err
	}
	return nil
}

func (s *ServiceCode) RemoveRegistercode(rawHex model.RegistercodeRawHex) error {
	if err := s.RepoRegistercode.Remove(rawHex); err != nil {
		return err
	}
	return nil
}
