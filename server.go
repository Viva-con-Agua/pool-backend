package main

import (
	"pool-backend/dao"
	"pool-backend/handlers/admin"
	"pool-backend/handlers/key"
	"pool-backend/handlers/token"

	"github.com/Viva-con-Agua/vcago"
)

func main() {
	e := vcago.NewServer()
	dao.InitialDatabase()
	dao.InitialNats()
	dao.InitialIDjango()
	dao.FixDatabase()
	dao.UpdateDatabase()
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
	token.Active.Routes(tokenUser.Group("/active"))
	token.NVM.Routes(tokenUser.Group("/nvm"))
	token.Address.Routes(tokenUser.Group("/address"))
	token.Avatar.Routes(tokenUser.Group("/avatar"))
	token.Newsletter.Routes(tokenUser.Group("/newsletter"))
	token.User.Routes(tokenUser)
	//crew routes
	crews := api.Group("/crews")
	token.Crew.Routes(crews)

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

	key.Import.Routes(api.Group("/import"))

	admin.Crew.Routes(e.Group("/admin/crews"))
	admin.Role.Routes(e.Group("/admin/users/role"))
	admin.User.Routes(e.Group("/admin/users"))

	//server
	e.Run()
}
