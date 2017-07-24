package auth

import (
	"github.com/labstack/echo"
	"net/http"
	"net/url"
)

const (
	GH_AUTHORIZE = "http://github.com/login/oauth/authorize"
)

func signIn(c echo.Context) error {
	authorizeUrl := url.Parse(GH_AUTHORIZE)
	q := authorizeUrl.Query()
	q.Add("client_id", "")
	q.Add("redirect_uri", "http://localhost:8000/oauth/github/")
	c.Redirect(http.StatusTemporaryRedirect, GH_AUTHORIZE)
	return nil
}
