package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	GH_AUTHORIZE             = "https://github.com/login/oauth/authorize"
	GH_TOKEN                 = "https://github.com/login/oauth/access_token"
	CONFIG_GH_CLIENT_ID      = "ghClientId"
	CONFIG_GH_CLIENT_SECRET  = "ghClientSecret"
	CONFIG_GH_REDIRECT       = "ghOauthRedirectUri"
	CONFIG_OAUTH_LANDING_URI = "oauthLandingRedirectUri"
	CONFIG_JWT_SECRET        = "jwtSecret"
)

// Challenge will redirect a user to authenticate with Github
func Challenge(c echo.Context) error {
	authorizeUrl, err := url.Parse(GH_AUTHORIZE)
	if err != nil {
		return err
	}

	q := authorizeUrl.Query()
	q.Add("client_id", viper.GetString(CONFIG_GH_CLIENT_ID))
	q.Add("redirect_uri", viper.GetString(CONFIG_GH_REDIRECT))
	q.Add("state", "42") // this should be replaced with a legit state variable
	authorizeUrl.RawQuery = q.Encode()
	c.Redirect(http.StatusTemporaryRedirect, authorizeUrl.String())
	return nil
}

func FetchCode(c echo.Context) error {
	ghError := c.QueryParam("error")
	if ghError != "" {
		return fmt.Errorf("%s: %s", ghError, c.QueryParam("error-description"))
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if state != "42" {
		return errors.New("You are a bad actor and are not allowed to access our site.")
	}

	values := url.Values{}
	values.Set("client_id", viper.GetString(CONFIG_GH_CLIENT_ID))
	values.Set("client_secret", viper.GetString(CONFIG_GH_CLIENT_SECRET))
	values.Set("code", code)
	values.Set("redirect_uri", viper.GetString(CONFIG_GH_REDIRECT))
	values.Set("state", state)
	resp, err := http.PostForm(GH_TOKEN, values)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	strBody := string(body)
	respValues, err := url.ParseQuery(strBody)
	ghError = respValues.Get("error")
	if ghError != "" {
		return fmt.Errorf("%s: %s", ghError, respValues.Get("error-description"))
	}

	landing, err := url.Parse(viper.GetString(CONFIG_OAUTH_LANDING_URI))
	if err != nil {
		return err
	}

	q := landing.Query()
	accessToken := respValues.Get("access_token")
	user, err := getGhUser(accessToken)
	if err != nil {
		return err
	}

	jwtToken, err := NewJwtToken(c, accessToken, user)
	if err != nil {
		return err
	}

	q.Add("token", jwtToken)
	landing.RawQuery = q.Encode()
	c.Redirect(http.StatusTemporaryRedirect, landing.String())

	return nil
}

func getGhUser(accessToken string) (*github.User, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// list all repositories for the authenticated user
	me, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return me, nil
}

func NewJwtToken(c echo.Context, accessToken string, user *github.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["login"] = user.Login
	claims["email"] = user.Email
	claims["ghToken"] = accessToken
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(viper.GetString(CONFIG_JWT_SECRET)))
	if err != nil {
		return "", err
	}
	return t, nil
}
