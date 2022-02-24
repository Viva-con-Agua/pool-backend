package handlers

type RoleUpdateRequest struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
}

/*
func RoleUpdateSet(c echo.Context) (err error) {
	ctx := c.Request().Context()
	body := new(RoleUpdateRequest)
	if vcago.BindAndValidate(c, body); err != nil {
		return
	}
	userReq := new(vcapool.User)
	if userReq, err = vcapool.AccessCookieUser(c); err != nil {
		return
	}
	var role *vcago.Role
	if role, err = vcago.Ge; err != nil {
		return
	}
	if !userReq.Roles.CheckRoot(role) {
		return vcago.NewValidationError("no permission for set this role")
	}
	user := new(dao.User)
	if err = user.Get(ctx, bson.M{"_id": body.UserID}); err != nil {
		return
	}
	user.Roles.Append(role)
	if err = user.Update(ctx); err != nil {
		return
	}
	return c.JSON(vcago.NewResponse("user.roles", user).Executed())
}*/
