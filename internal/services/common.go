package services

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"strings"
)

func HandleLogin(c echo.Context, oauthConf *oauth2.Config, oauthStateString string) error {
	// Build the OAuth URL
	URL := oauthConf.Endpoint.AuthURL + "?" + url.Values{
		"client_id":     {oauthConf.ClientID},
		"scope":         {strings.Join(oauthConf.Scopes, " ")},
		"redirect_uri":  {oauthConf.RedirectURL},
		"response_type": {"code"},
		"state":         {oauthStateString},
	}.Encode()

	// Log information
	c.Logger().Info("Login URL:", URL)

	// Redirect user to the OAuth provider
	return c.Redirect(http.StatusTemporaryRedirect, URL)
}
