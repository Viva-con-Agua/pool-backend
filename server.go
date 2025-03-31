package main

import (
	"pool-backend/dao"
	_ "pool-backend/docs"
	"pool-backend/handlers/admin"
	"pool-backend/handlers/key"
	"pool-backend/handlers/token"

	"github.com/Viva-con-Agua/vcago"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Pool API documentation
// @version 3.0.0
// @host pool3-api.vivaconagua.org
// @schemes https

// @BasePath /v1

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://hydra.vivaconagua.org/oauth/token
// @authorizationurl https://hydra.vivaconagua.org/oauth/authorize

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Bearer key

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name access_cookie
// @description

func main() {
	e := vcago.NewServer()
	dao.InitialDatabase()
	dao.InitialNats()
	dao.InitialIDjango()
	dao.InitialTestLogin()
	dao.FixDatabase()
	dao.UpdateDatabase()
	dao.UpdateTicker()
	//dao.ReloadDatabase()
	//login routes
	api := e.Group("/v1")

	token.Assets.Routes(api.Group("/assets"))

	token.Login.Routes(api.Group("/auth"))
	//user routes
	tokenUser := api.Group("/users")
	token.User.Routes(tokenUser)
	token.Profile.Routes(tokenUser.Group("/profile"))
	token.UserCrew.Routes(tokenUser.Group("/crew"))
	token.Role.Routes(tokenUser.Group("/role"))
	token.RoleHistory.Routes(tokenUser.Group("/role_history"))
	token.Active.Routes(tokenUser.Group("/active"))
	token.NVM.Routes(tokenUser.Group("/nvm"))
	token.Address.Routes(tokenUser.Group("/address"))
	token.Avatar.Routes(tokenUser.Group("/avatar"))
	token.Newsletter.Routes(tokenUser.Group("/newsletter"))

	token.User.Routes(tokenUser)
	//crew routes
	crews := api.Group("/crews")
	token.Crew.Routes(crews)

	//crew routes
	organisations := api.Group("/organisations")
	token.Organisation.Routes(organisations)

	mails := api.Group("/mails")
	token.Mailbox.Routes(mails.Group("/mailbox"))
	token.Message.Routes(mails.Group("/message"))

	events := api.Group("/events")
	token.Event.Routes(events.Group("/event"))
	token.Artist.Routes(events.Group("/artist"))
	token.Organizer.Routes(events.Group("/organizer"))
	token.Participation.Routes(events.Group("/participation"))

	finances := api.Group("/finances")
	token.Source.Routes(finances.Group("/source"))
	token.Taking.Routes(finances.Group("/taking"))
	token.Deposit.Routes(finances.Group("/deposit"))
	token.ReceiptFile.Routes(finances.Group("/receipt"))

	key.Import.Routes(api.Group("/import"))

	if dao.TestLogin {
		admin.Crew.Routes(e.Group("/admin/crews"))
		admin.Role.Routes(e.Group("/admin/users/role"))
		admin.User.Routes(e.Group("/admin/users"))
	}
	api.GET("/docu/*", echoSwagger.WrapHandler)
	//server
	e.Run()
}
