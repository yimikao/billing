package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	applogger "github.com/fluidcoins/log"
	"github.com/yimikao/billing/core/oauth"
	"golang.org/x/oauth2"
)

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)

	state := base64.URLEncoding.EncodeToString(b)

	cookie := http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}

	http.SetCookie(w, &cookie)

	// return cookie value
	return state
}

type LoginHandler struct {
	cfg    *oauth2.Config
	logger applogger.Entry
}

func NewLoginHandler(cfg *oauth2.Config, logger applogger.Entry) *LoginHandler {
	return &LoginHandler{
		cfg:    cfg,
		logger: logger,
	}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {

	// add cookie to response and also return it's value
	// 'll be validated that it matches with the state query parameter on redirect callback
	oauthState := generateStateOauthCookie(w)

	// redirect user to googleauth url
	u := h.cfg.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)

}

type CallbackHandler struct {
	client *oauth.GoogleOauthClient
	logger applogger.Entry
}

func NewCallbackHandler(client *oauth.GoogleOauthClient, logger applogger.Entry) *CallbackHandler {
	return &CallbackHandler{
		client: client,
		logger: logger,
	}
}

func (h *CallbackHandler) Callback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthStateCookie, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthStateCookie.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		log.Println("auth code not supplied")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := h.client.GetUserDataFromGoogle(code)
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var ur = new(UserResponse)
	if err := json.Unmarshal(data, ur); err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "UserInfo: %s\n %s\n", ur.Email, ur.Picture)

	// t, _ := template.ParseFiles("templates/success.html")
	// t.Execute(w, ur)
	// fmt.Fprintf(w, "UserInfo: %s\n", data)
}
