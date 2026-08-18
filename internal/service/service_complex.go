package service

type ServiceComplex struct {
	ServiceUser  *ServiceUser
	ServiceCode  *ServiceCode
	ServiceAdmin *ServiceAdmin
}

func NewServiceComplex(
	pSvcUser *ServiceUser,
	pSvcCode *ServiceCode,
	pSvcAdmin *ServiceAdmin,
) *ServiceComplex {
	return &ServiceComplex{
		ServiceUser:  pSvcUser,
		ServiceCode:  pSvcCode,
		ServiceAdmin: pSvcAdmin,
	}
}
