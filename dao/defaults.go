package dao

import (
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
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

func (i *Check) ASP(token *models.AccessToken) *Check {
	if !token.PoolRoles.Validate("finance;network;education;") {
		i.Error = vcago.NewPermissionDenied("permission_denied", nil)
	}
	return i
}
