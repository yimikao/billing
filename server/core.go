package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fluidcoins/log"
	"github.com/go-chi/render"
	"github.com/yimikao/billing/core"
	"golang.org/x/oauth2"
)

type LoginHandler struct {
	OauthConfig *oauth2.Config
	OauthState  string
	logger      log.Entry
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		OauthConfig: core.OauthConfig,
		OauthState:  core.OauthState,
	}
}

func NewCallbackHandler() *CallbackHandler {
	return &CallbackHandler{
		OauthConfig: core.OauthConfig,
		OauthState:  core.OauthState,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	URL, err := url.Parse(h.OauthConfig.Endpoint.AuthURL)

	if err != nil {
		h.logger.WithError(err).Error(err.Error())
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request body"))
		return
	}

	params := url.Values{}
	params.Add("client_id", h.OauthConfig.ClientID)
	params.Add("scope", strings.Join(h.OauthConfig.Scopes, " "))
	params.Add("redirect_uri", h.OauthConfig.RedirectURL)
	params.Add("response_type", "code")
	params.Add("state", h.OauthState)

	URL.RawQuery = params.Encode()

	url := URL.String()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

type CallbackHandler struct {
	OauthConfig *oauth2.Config
	OauthState  string
	Logger      log.Entry
}

func (h *CallbackHandler) Callback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != h.OauthState {
		// h.logger.Error("invalid oauth state")
		// _ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid oauth state"))
		fmt.Println("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	if code == "" {

		// _ = render.Render(w, r, newAPIError(http.StatusBadRequest, "code not found to provide token"))
		fmt.Println("code to provide token not found")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			// _ = render.Render(w, r, newAPIError(http.StatusBadRequest, "user has denied perm"))
			fmt.Println("user has denied perm")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}

	} else {
		token, err := h.OauthConfig.Exchange(context.Background(), code)
		if err != nil {
			fmt.Println("oauth exchange failed" + err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
		if err != nil {
			// logger.Log.Error("Get: " + err.Error() + "\n")
			fmt.Println("getting body error" + err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		defer resp.Body.Close()

		response, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("reading body err" + err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Println(response)
		_ = render.Render(w, r, newAPIStatus(http.StatusOK, "login successful"))

	}

}
