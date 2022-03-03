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
	e.Use(vcago.Logger.Init())
	e.Use(vcago.CORS.Init())

	//error
	e.HTTPErrorHandler = vcago.HTTPErrorHandler
	e.Validator = vcago.JSONValidator
	login := e.Group("/auth")
	login.POST("/callback", handlers.CallbackHandler)
	login.GET("/refresh", handlers.RefreshHandler, vcapool.RefreshCookieConfig())
	users := e.Group("/users")
	users.GET("", handlers.ListUser, vcapool.AccessCookieConfig())

	crewUser := users.Group("/crew")
	crewUser.POST("", handlers.CreateUserCrew, vcapool.AccessCookieConfig())

	roles := users.Group("/role")
	roles.POST("", handlers.RoleCreate, vcapool.AccessCookieConfig())
	roles.DELETE("", handlers.RoleDelete, vcapool.AccessCookieConfig())

	activeUser := users.Group("/active")
	activeUser.GET("/request", handlers.RequestUserActive, vcapool.AccessCookieConfig())
	activeUser.POST("/confirm", handlers.ConfirmUserActive, vcapool.AccessCookieConfig())
	activeUser.POST("/reject", handlers.RejectUserActive, vcapool.AccessCookieConfig())

	address := users.Group("/address")
	address.POST("", handlers.CreateAddress, vcapool.AccessCookieConfig())
	address.PUT("", handlers.UpdateAddress, vcapool.AccessCookieConfig())

	crews := e.Group("/crews")
	crews.POST("", handlers.CreateCrew)
	crews.GET("", handlers.ListCrew)
	crews.GET("/:id", handlers.GetCrew)
	crews.PUT("", handlers.UpdateCrew)
	crews.DELETE("", handlers.DeleteCrew)

	test := e.Group("/test/users")
	test.GET("", handlers.ListUser)
	//server
	appPort := vcago.Config.GetEnvString("APP_PORT", "n", "1323")

	e.Logger.Fatal(e.Start(":" + appPort))

}
