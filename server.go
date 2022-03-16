package main

import (
	"pool-user/handlers"

	"github.com/Viva-con-Agua/vcago"
	"github.com/Viva-con-Agua/vcapool"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Debug = false
	// Middleware
	e.Use(vcago.Logger.Init("pool-user"))
	e.Use(vcago.CORS.Init())
	vcago.Nats.Connect()
	//error
	e.HTTPErrorHandler = vcago.HTTPErrorHandler
	e.Validator = vcago.JSONValidator
	login := e.Group("/auth")
	login.POST("/callback", handlers.CallbackHandler)
	login.GET("/refresh", handlers.RefreshHandler, vcapool.RefreshCookieConfig())
	login.GET("/logout", handlers.LogoutHandler, vcapool.AccessCookieConfig())

	users := e.Group("/users")
	users.GET("", handlers.ListUser, vcapool.AccessCookieConfig())

	profile := users.Group("/profile")
	profile.POST("", handlers.ProfileCreate, vcapool.AccessCookieConfig())
	profile.PUT("", handlers.ProfileUpdate, vcapool.AccessCookieConfig())

	crewUser := users.Group("/crew")
	crewUser.POST("", handlers.CreateUserCrew, vcapool.AccessCookieConfig())
	crewUser.PUT("", handlers.UpdateUserCrew, vcapool.AccessCookieConfig())
	crewUser.DELETE("", handlers.DeleteUserCrew, vcapool.AccessCookieConfig())

	roles := users.Group("/role")
	roles.POST("", handlers.RoleCreate, vcapool.AccessCookieConfig())
	roles.DELETE("", handlers.RoleDelete, vcapool.AccessCookieConfig())

	activeUser := users.Group("/active")
	activeUser.GET("/request", handlers.RequestUserActive, vcapool.AccessCookieConfig())
	activeUser.POST("/confirm", handlers.ConfirmUserActive, vcapool.AccessCookieConfig())
	activeUser.POST("/reject", handlers.RejectUserActive, vcapool.AccessCookieConfig())
	activeUser.GET("/withdraw", handlers.WithdrawUserActive, vcapool.AccessCookieConfig())

	nvmUser := users.Group("/nvm")
	nvmUser.GET("/confirm", handlers.ConfirmUserNVM, vcapool.AccessCookieConfig())
	nvmUser.POST("/reject", handlers.RejectUserNVM, vcapool.AccessCookieConfig())
	nvmUser.GET("/withdraw", handlers.WithdrawUserNVM, vcapool.AccessCookieConfig())

	address := users.Group("/address")
	address.POST("", handlers.CreateAddress, vcapool.AccessCookieConfig())
	address.PUT("", handlers.UpdateAddress, vcapool.AccessCookieConfig())

	avatar := users.Group("/avatar")
	avatar.POST("", handlers.CreateAvatar, vcapool.AccessCookieConfig())
	avatar.DELETE("", handlers.DeleteAddress, vcapool.AccessCookieConfig())

	crews := e.Group("/crews")
	crews.POST("", handlers.CreateCrew, vcago.AccessCookieConfig())
	crews.GET("", handlers.ListCrew, vcago.AccessCookieConfig())
	crews.GET("/:id", handlers.GetCrew, vcago.AccessCookieConfig())
	crews.PUT("", handlers.UpdateCrew, vcago.AccessCookieConfig())
	crews.DELETE("", handlers.DeleteCrew, vcago.AccessCookieConfig())

	test := e.Group("/test/users")
	test.GET("", handlers.ListUser)
	//server
	appPort := vcago.Config.GetEnvString("APP_PORT", "n", "1323")

	e.Logger.Fatal(e.Start(":" + appPort))

}
