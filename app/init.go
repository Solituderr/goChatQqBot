package app

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

var e *echo.Echo

func InitWebFramework() {
	e = echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	addRoutes()

	logrus.Info("echo framework initialized")
}

func StartServer() {
	e.Logger.Fatal(e.Start(":5703"))
}
