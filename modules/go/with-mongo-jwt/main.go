package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ruthv1k/flock/modules/go/with-mongo-jwt/controllers"
	mongoconnect "github.com/ruthv1k/flock/modules/go/with-mongo-jwt/database"
)

var API_VERSION = "/api/v1"

func main() {
	e := echo.New()

	// public route groups
	authRoutes := e.Group(API_VERSION + "/auth")

	// public routes
	authRoutes.POST("/register", controllers.Register)
	authRoutes.POST("/login", controllers.Login)

	e.Logger.Fatal(e.Start(":5000"))
	defer mongoconnect.DisconnectDb()
}
