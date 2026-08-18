package service

import (
	abstract "github.com/atomreforge/daizy-night-server/internal/abstract/interface"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

var _ abstract.InterfaceServiceCode = (*ServiceCode)(nil)

type ServiceCode struct {
	repoRegistercode abstract.InterfaceRepoRegistercode
}

func NewServiceCode(repoRegcode abstract.InterfaceRepoRegistercode) *ServiceCode {
	return &ServiceCode{
		repoRegistercode: repoRegcode,
	}
}

func (s *ServiceCode) RecordNewRegistercode(registercodeRecord *model.RegistercodeRecord) error {
	if err := s.repoRegistercode.RecordNewRegistercode(registercodeRecord); err != nil {
		return err
	}
	return nil
}

func (s *ServiceCode) RemoveRegistercode(rawHex model.RegistercodeRawHex) error {
	if err := s.repoRegistercode.Remove(rawHex); err != nil {
		return err
	}
	return nil
}
