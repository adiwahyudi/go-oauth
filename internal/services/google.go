package services

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
)

var (
	oauthConfGl = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:9090/callback-google",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateStringGl = ""
)

func InitializeOAuthGoogle() {
	oauthConfGl.ClientID = viper.GetString("google.clientID")
	oauthConfGl.ClientSecret = viper.GetString("google.clientSecret")
	oauthStateStringGl = viper.GetString("oauthStateString")
}

func HandleGoogleLogin(c echo.Context) error {
	return HandleLogin(c, oauthConfGl, oauthStateStringGl)
}
func HandleGoogleCallback(c echo.Context) error {
	state := c.FormValue("state")
	if state != oauthStateStringGl {
		return c.String(http.StatusInternalServerError, "States don't Match!!")
	}

	code := c.FormValue("code")
	if code == "" {
		return c.String(http.StatusInternalServerError, "Code not found!")
	}

	token, err := oauthConfGl.Exchange(context.Background(), code)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "JSON Parsing Failed")
	}

	return c.String(http.StatusOK, string(userData))
}
