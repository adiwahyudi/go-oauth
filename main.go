package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"go-oauth/internal/configs"
	"go-oauth/internal/services"
)

func main() {

	e := echo.New()
	configs.InitializeViper()
	services.InitializeOAuthGoogle()
	services.InitializeOAuthFacebook()
	e.GET("/", func(c echo.Context) error {
		return c.File("assets/index.html")
	})

	// OAuth
	e.GET("/login-google", services.HandleGoogleLogin)
	e.GET("/callback-google", services.HandleGoogleCallback)

	e.GET("/login-facebook", services.HandleFacebookLogin)
	e.GET("/callback-facebook", services.HandleCallbackFacebook)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", viper.GetString("port"))))
}
