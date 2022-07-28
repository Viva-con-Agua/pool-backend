package main

import (
	"pool-user/handlers/admin"
	"pool-user/handlers/key"
	"pool-user/handlers/token"

	"github.com/Viva-con-Agua/vcago"
)

func main() {
	e := vcago.NewEchoServer("pool-user")
	vcago.Nats.Connect()
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

	key.Crew.Routes(e.Group("/apikey/crews"))

	admin.Crew.Routes(e.Group("/admin/crews"))
	admin.Role.Routes(e.Group("/admin/users/role"))
	admin.User.Routes(e.Group("/admin/users"))

	//server
	port := vcago.Config.GetEnvString("APP_PORT", "n", "1323")
	e.Logger.Fatal(e.Start(":" + port))
}
