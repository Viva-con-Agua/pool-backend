package token

import (
	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
)

var refreshCookie = vcago.RefreshCookieMiddleware()
var accessCookie = vcago.AccessCookieMiddleware(&vcapool.AccessToken{})
