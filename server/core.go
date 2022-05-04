package server

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	applogger "github.com/fluidcoins/log"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"
	"github.com/yimikao/billing"
	"github.com/yimikao/billing/core/oauth"
	redisDB "github.com/yimikao/billing/database/redis"
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
	client   *oauth.GoogleOauthClient
	userRepo billing.UserRepository
	logger   applogger.Entry
}

func NewCallbackHandler(client *oauth.GoogleOauthClient, userRepo billing.UserRepository, logger applogger.Entry) *CallbackHandler {
	return &CallbackHandler{
		client:   client,
		userRepo: userRepo,
		logger:   logger,
	}
}

func (h *CallbackHandler) Callback(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		//
		h.logger.WithError(err).Error("request body malformed")
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request url queries"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthStateCookie, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthStateCookie.Value {

		h.logger.Error("invalid oauth google state")
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request url queries. invalid oauth gogle state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	code := r.FormValue("code")
	if code == "" {

		h.logger.Error("auth code not supplied")
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request url queries. auth code not supplied"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	data, err := h.client.GetUserDataFromGoogle(code)
	if err != nil {

		h.logger.WithError(err).Error("could not get user data from google with supplied code")
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request url queries"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	var ur = new(UserResponse)
	if err := json.Unmarshal(data, ur); err != nil {

		h.logger.Info("user data gotten successfully")

	}

	// fmt.Fprintf(w, "UserInfo: %s\n %s\n", ur.Email, ur.Picture)

	u, err := h.userRepo.CheckAlreadyRegistered(ur.Email)

	if err != nil {
		// go to registration flow.
		return
	}

	_ = u
	// load user account data
}

type UserRegistrationHandler struct {
	userRepo billing.UserRepository
	logger   applogger.Entry
	// redisClient *redisDB.Client
	cache   redisDB.Cache
	context context.Context
}

func NewUserRegistrationHandler(ur billing.UserRepository, l applogger.Entry, cache redisDB.Cache, ctx context.Context) *UserRegistrationHandler {
	return &UserRegistrationHandler{
		userRepo: ur,
		logger:   l,
		// redisClient: rc,
		cache:   cache,
		context: ctx,
	}
}

type userRegistrationRequest struct {
	// will be fetched automatically from google account
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`

	// to be entered by user
	Tag             string `json:"tag"`
	TransactionCode string `json:"transaction_code"`
	GenericRequest
}

func (req *userRegistrationRequest) validate() error {

	if strings.TrimSpace(req.Tag) == "" {
		return errors.New("please provide a valid tag")
	}

	if len(strings.TrimSpace(req.Tag)) < 4 {
		return errors.New("tag cannot be less than four characters")
	}

	if len(strings.TrimSpace(req.Tag)) > 15 {
		return errors.New("tag cannot be more than 15 characters")
	}

	if strings.TrimSpace(req.TransactionCode) == "" {
		return errors.New("please provide transaction code")
	}

	if _, err := strconv.Atoi(strings.TrimSpace(req.TransactionCode)); err != nil {
		return errors.New("transaction code must be digits alone")
	}

	if len(strings.TrimSpace(req.TransactionCode)) < 4 {
		return errors.New("transaction code cannot be less than four digits")
	}

	if len(strings.TrimSpace(req.TransactionCode)) > 4 {
		return errors.New("transaction code cannot be more than four digits")
	}

	return nil

}

func (h *UserRegistrationHandler) registerUser(w http.ResponseWriter, r *http.Request) {

	var req = new(userRegistrationRequest)

	if err := render.Bind(r, req); err != nil {
		h.logger.WithError(err).Error("request body malformed")
		_ = render.Render(w, r, errInvalidRequestBody)
		return
	}

	if err := req.validate(); err != nil {
		h.logger.WithError(err).Error("request body malformed")
		_ = render.Render(w, r, errInvalidRequestBody)
		return
	}

	us := &billing.User{
		Tag:             req.Tag,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Email:           req.Email,
		TransactionCode: req.TransactionCode,
	}

	if err := h.userRepo.Create(us); err != nil {
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	b, _ := json.Marshal(us)
	_ = b

	if err := h.cache.StoreData(h.context, redisDB.Userdata, "hello"); err != nil {
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	if err := json.NewEncoder(w).Encode(us); err != nil {
		log.Fatal(err)
	}

	// _ = render.Render(w, r, newAPIStatus(http.StatusOK, "registration sucessfull"))
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)

}

type HomepageHandler struct {
	// userData *billing.User
	logger applogger.Entry
	// redisClient *redisDB.Client
	cache   redisDB.Cache
	context context.Context
}

func NewHomepageHandler(logger applogger.Entry, cache redisDB.Cache, ctx context.Context) *HomepageHandler {
	return &HomepageHandler{
		// userData: userData,
		logger: logger,
		// redisClient: redisClient,
		cache:   cache,
		context: ctx,
	}
}

func (h *HomepageHandler) displayUserData(w http.ResponseWriter, r *http.Request) {

	userData, err := h.cache.GetData(h.context, redisDB.Userdata)

	if err != nil {
		if err == redis.Nil {
			h.logger.WithError(err).Error(err.Error())
			_ = render.Render(w, r, newAPIError(http.StatusNotFound, "key does not exist"))
			return
		}

		h.logger.WithError(err).Error(err.Error())
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError, "could not get user data"))
		return

	}

	u := new(billing.User)

	if err := json.Unmarshal([]byte(userData), u); err != nil {
		h.logger.WithError(err).Error(err.Error())
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError, "could not load user data"))
		return
	}

}
