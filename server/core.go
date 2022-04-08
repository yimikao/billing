package server

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/fluidcoins/log"
	"github.com/go-chi/render"
	"golang.org/x/oauth2"
)

type LoginHandler struct {
	oauthConfig *oauth2.Config
	oauthState  string
	logger      log.Entry
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	URL, err := url.Parse(h.oauthConfig.Endpoint.AuthURL)

	if err != nil {
		h.logger.WithError(err).Error(err.Error())
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request body"))
		return
	}

	params := url.Values{}
	params.Add("client_id", h.oauthConfig.ClientID)
	params.Add("scope", strings.Join(h.oauthConfig.Scopes, " "))
	params.Add("redirect_uri", h.oauthConfig.RedirectURL)
	params.Add("response_type", "code")
	params.Add("state", h.oauthState)

	URL.RawQuery = params.Encode()

	url := URL.String()

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

type CallbackHandler struct {
	oauthConfig *oauth2.Config
	oauthState  string
	logger      log.Entry
}

func (h *CallbackHandler) Callback(w http.ResponseWriter, r *http.Request) {

	state := r.FormValue("state")
	if state != h.oauthState {
		// h.logger.WithError(err).Error(err.Error())
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid oauth state"))
		return
	}

}
