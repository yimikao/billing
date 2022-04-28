package oauth

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/yimikao/billing/core"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOauthClient struct {
	cfg               *oauth2.Config
	oauthGoogleUrlAPI string
}

func NewGoogleOauthClient(cfg *oauth2.Config) *GoogleOauthClient {
	return &GoogleOauthClient{
		cfg:               cfg,
		oauthGoogleUrlAPI: "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
	}
}

func NewGoogleOauthConfig(cfg core.Config) *oauth2.Config {

	return &oauth2.Config{
		ClientID:     cfg.Oauth.OauthClientID,
		ClientSecret: cfg.Oauth.OauthClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  cfg.Oauth.OauthCallbackURL,
		Scopes:       cfg.Oauth.OauthScopes,
	}

}

func (c *GoogleOauthClient) GetUserDataFromGoogle(code string) ([]byte, error) {
	// use auth-code to get token and get user info

	token, err := c.cfg.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	res, err := http.Get(c.oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err.Error())
	}

	return bts, nil

}
