package services

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"io"
	"log"
	"net/http"
)

var (
	oauthConfFb = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9090/callback-facebook",
		Scopes: []string{
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}
	oauthStateStringFb = ""
)

func InitializeOAuthFacebook() {
	oauthConfFb.ClientID = viper.GetString("facebook.clientID")
	oauthConfFb.ClientSecret = viper.GetString("facebook.clientSecret")
	oauthStateStringFb = viper.GetString("oauthStateString")
}

func HandleFacebookLogin(c echo.Context) error {
	return HandleLogin(c, oauthConfFb, oauthStateStringFb)
}

func HandleCallbackFacebook(c echo.Context) error {
	state := c.FormValue("state")
	if state != oauthStateStringFb {
		return c.String(http.StatusInternalServerError, "States don't Match!!")
	}

	code := c.FormValue("code")
	if code == "" {
		return c.String(http.StatusInternalServerError, "Code not found!")
	}

	token, err := oauthConfFb.Exchange(context.Background(), code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
	}

	client := oauthConfFb.Client(context.Background(), token)
	resp, err := client.Get("https://graph.facebook.com/v19.0/me?fields=id,email,name")
	if err != nil {
		return c.String(http.StatusInternalServerError, "User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "JSON Parsing Failed")
	}

	log.Println(string(userData))

	return c.String(http.StatusOK, string(userData))

}
