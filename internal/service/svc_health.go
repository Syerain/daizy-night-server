package service

import "github.com/atomreforge/daizy-night-server/internal/dbware"

type ServiceHealth struct {
	pDB *dbware.ProviderDB
}

func NewServiceHealth(pDB *dbware.ProviderDB) *ServiceHealth {
	return &ServiceHealth{
		pDB: pDB,
	}
}

func (p *ServiceHealth) HealthCheckDb() (bool, error) {
	err := p.pDB.Check()
	if err != nil {
		return false, err
	}
	return true, nil
}
