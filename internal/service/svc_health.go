package service

import "github.com/atomreforge/daizy-night-server/internal/dbware"

type ProviderHealth struct {
	pDB *dbware.ProviderDB
}

func (p *ProviderHealth) HealthCheckDb() (bool, error) {
	err := p.pDB.Check()
	if err != nil {
		return false, err
	}
	return true, nil
}
