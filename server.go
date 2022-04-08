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

	loginHandler := handlers.NewLoginHandler()
	login := e.Group("/auth")
	login.Use(loginHandler.Context)
	login.POST("/callback", loginHandler.Callback)
	if vcago.Config.GetEnvBool("API_TEST_LOGIN", "n", false) {
		login.POST("/testlogin", loginHandler.LoginAPI)
	}
	login.GET("/refresh", loginHandler.Refresh, vcapool.RefreshCookieConfig())
	login.GET("/logout", loginHandler.Logout, vcago.AccessCookieMiddleware(&vcapool.AccessToken{}))

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

	addressHandler := handlers.NewAddressHandler()
	address := users.Group("/address")
	address.Use(addressHandler.Context)
	address.POST("", addressHandler.Create, vcapool.AccessCookieConfig())
	address.PUT("", addressHandler.Update, vcapool.AccessCookieConfig())
	address.GET("/:id", addressHandler.Get, vcapool.AccessCookieConfig())
	address.DELETE("/:id", addressHandler.Delete, vcapool.AccessCookieConfig())

	avatarHandler := handlers.NewAvatarHandler()
	avatar := users.Group("/avatar")
	avatar.Use(avatarHandler.Context)
	avatar.POST("", avatarHandler.Create, vcapool.AccessCookieConfig())
	avatar.PUT("", avatarHandler.Update, vcapool.AccessCookieConfig())
	avatar.DELETE("/:id", avatarHandler.Delete, vcapool.AccessCookieConfig())

	crews := e.Group("/crews")
	crews.POST("", handlers.CrewCreate, vcapool.AccessCookieConfig())
	crews.GET("", handlers.CrewList)
	crews.GET("/:id", handlers.CrewGet)
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
