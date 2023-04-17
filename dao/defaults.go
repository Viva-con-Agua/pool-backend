package dao

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
)

type Check struct {
	Error error
}

func NewCheck() *Check {
	return &Check{}
}

func (i *Check) Return() error {
	return i.Error
}

func (i *Check) ASP(token *vcapool.AccessToken) *Check {
	if !token.PoolRoles.Validate("finance;network;education;") {
		i.Error = vcago.NewBadRequest("permission", "permission_denied", nil)
	}
	return i
}
