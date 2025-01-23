package token

import (
	"pool-backend/models"

	"github.com/Viva-con-Agua/vcago"
)

var refreshCookie = vcago.RefreshCookieMiddleware()
var accessCookie = vcago.AccessCookieMiddleware(&models.AccessToken{})
