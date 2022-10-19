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
	//dao.ReloadDatabase()
	//login routes
	token.Login.Routes(e.Group("/auth"))

	//user routes
	tokenUser := e.Group("/users")
	token.User.Routes(tokenUser)
	token.Profile.Routes(tokenUser.Group("/profile"))
	token.UserCrew.Routes(tokenUser.Group("/crew"))
	token.Role.Routes(tokenUser.Group("/role"))
	token.Active.Routes(tokenUser.Group("/active"))
	token.NVM.Routes(tokenUser.Group("/nvm"))
	token.Address.Routes(tokenUser.Group("/address"))
	token.Avatar.Routes(tokenUser.Group("/avatar"))
	token.User.Routes(tokenUser)
	//crew routes
	crews := e.Group("/crews")
	token.Crew.Routes(crews)

	mails := e.Group("/mails")
	token.Mailbox.Routes(mails.Group("/mailbox"))
	token.Message.Routes(mails.Group("/message"))

	events := e.Group("/events")
	token.Event.Routes(events.Group("/event"))
	token.Artist.Routes(events.Group("/artist"))
	token.Organizer.Routes(events.Group("/organizer"))
	token.Participation.Routes(events.Group("/participation"))

	key.Crew.Routes(e.Group("/apikey/crews"))

	admin.Crew.Routes(e.Group("/admin/crews"))
	admin.Role.Routes(e.Group("/admin/users/role"))
	admin.User.Routes(e.Group("/admin/users"))

	//server
	port := vcago.Settings.String("APP_PORT", "n", "1323")
	e.Logger.Fatal(e.Start(":" + port))
}
