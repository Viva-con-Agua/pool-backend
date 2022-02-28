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
	api := e.Group("/users")

	api.POST("/crew", handlers.CreateUserCrew, vcapool.AccessCookieConfig())

	api.GET("/active/request", handlers.CreateUserActive, vcapool.AccessCookieConfig())

	api.Use(vcapool.AccessCookieConfig())
	address := api.Group("/address")
	address.POST("", handlers.CreateAddress)
	address.PUT("", handlers.UpdateAddress)

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
