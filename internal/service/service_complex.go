package service

type ServiceComplex struct {
	ServiceUser   *ServiceUser
	ServiceCode   *ServiceCode
	ServiceAdmin  *ServiceAdmin
	ServiceHealth *ServiceHealth
}

func NewServiceComplex(
	pSvcUser *ServiceUser,
	pSvcCode *ServiceCode,
	pSvcAdmin *ServiceAdmin,
	pSvcHealth *ServiceHealth,
) *ServiceComplex {
	return &ServiceComplex{
		ServiceUser:   pSvcUser,
		ServiceCode:   pSvcCode,
		ServiceAdmin:  pSvcAdmin,
		ServiceHealth: pSvcHealth,
	}
}
