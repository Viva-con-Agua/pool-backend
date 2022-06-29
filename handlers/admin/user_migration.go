package admin

/*
import (
	"pool-user/dao"

	"github.com/Viva-con-Agua/vcago"
	"github.com/labstack/echo/v4"
)

type UserMigrateHandler struct {
	vcago.Handler
}

func NewUserMigrationHandler() *UserMigrateHandler {
	handler := vcago.NewHandler("user_migration")
	return &UserMigrateHandler{
		*handler,
	}
}

func (i *UserMigrateHandler) Create(cc echo.Context) (err error) {
	c := cc.(vcago.Context)
	body := new(dao.UserMigrate)
	if err = c.BindAndValidate(body); err != nil {
		return
	}
	if err = body.MigrateUser(c.Ctx()); err != nil {
		return
	}
	return c.Created(body)
}
*/
