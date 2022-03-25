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
	users.GET("", handlers.UserList, vcapool.AccessCookieConfig())

	profile := users.Group("/profile")
	profile.POST("", handlers.ProfileCreate, vcapool.AccessCookieConfig())
	profile.PUT("", handlers.ProfileUpdate, vcapool.AccessCookieConfig())

	crewUser := users.Group("/crew")
	crewUser.POST("", handlers.UserCrewCreate, vcapool.AccessCookieConfig())
	crewUser.PUT("", handlers.UserCrewUpdate, vcapool.AccessCookieConfig())
	crewUser.DELETE("", handlers.UserCrewDelete, vcapool.AccessCookieConfig())

	roles := users.Group("/role")
	roles.POST("", handlers.RoleCreate, vcapool.AccessCookieConfig())
	roles.DELETE("", handlers.RoleDelete, vcapool.AccessCookieConfig())

	activeUser := users.Group("/active")
	activeUser.GET("/request", handlers.UserActiveRequest, vcapool.AccessCookieConfig())
	activeUser.POST("/confirm", handlers.UserActiveConfirm, vcapool.AccessCookieConfig())
	activeUser.POST("/reject", handlers.UserActiveReject, vcapool.AccessCookieConfig())
	activeUser.GET("/withdraw", handlers.UserActiveWithdraw, vcapool.AccessCookieConfig())

	nvmUser := users.Group("/nvm")
	nvmUser.GET("/confirm", handlers.UserNVMConfirm, vcapool.AccessCookieConfig())
	nvmUser.POST("/reject", handlers.UserNVMReject, vcapool.AccessCookieConfig())
	nvmUser.GET("/withdraw", handlers.UserNVMWithdraw, vcapool.AccessCookieConfig())

	address := users.Group("/address")
	address.POST("", handlers.AddressCreate, vcapool.AccessCookieConfig())
	address.PUT("", handlers.AddressUpdate, vcapool.AccessCookieConfig())
	address.GET("/:id", handlers.AddressGet, vcapool.AccessCookieConfig())
	address.DELETE("/:id", handlers.AddressDelete, vcapool.AccessCookieConfig())

	avatar := users.Group("/avatar")
	avatar.POST("", handlers.AvatarCreate, vcapool.AccessCookieConfig())
	avatar.DELETE("", handlers.AvatarDelete, vcapool.AccessCookieConfig())

	crews := e.Group("/crews")
	crews.POST("", handlers.CrewCreate, vcapool.AccessCookieConfig())
	crews.GET("", handlers.CrewList, vcapool.AccessCookieConfig())
	crews.GET("/:id", handlers.CrewGet, vcapool.AccessCookieConfig())
	crews.PUT("", handlers.CrewUpdate, vcapool.AccessCookieConfig())
	crews.DELETE("", handlers.CrewDelete, vcapool.AccessCookieConfig())

	admin := e.Group("/admin")
	adminUser := admin.Group("/users")
	adminUser.GET("", handlers.UserListAdmin)

	adminCrew := admin.Group("/crews")
	adminCrew.GET("", handlers.CrewListAdmin)

	//server
	port := vcago.Config.GetEnvString("APP_PORT", "n", "1323")

	e.Logger.Fatal(e.Start(":" + port))

}
