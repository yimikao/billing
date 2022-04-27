package server

import (
	"errors"
	"net/http"

	"github.com/fluidcoins/hermes/util"
	"github.com/fluidcoins/log"
	"github.com/go-chi/render"
	"github.com/yimikao/billing"
)

type GenericRequest struct{}

func (g GenericRequest) Bind(_ *http.Request) error { return nil }

type UserHandler struct {
	userRepo billing.UserRepository
	logger   log.Entry
}

func NewUserHandler(logger log.Entry, userRepo billing.UserRepository) *UserHandler {

	return &UserHandler{
		logger:   logger,
		userRepo: userRepo,
	}
}

type userCreateRequest struct {
	Tag       string        `json:"tag"`
	FirstName string        `json:"first_name"`
	LastName  string        `json:"last_name"`
	Password  string        `json:"password"`
	Email     billing.Email `json:"email"`
	GenericRequest
}

func (u *userCreateRequest) Validate() error {

	if util.IsStringEmpty(u.Email.String()) {
		return errors.New("please provide an email address")
	}

	return nil
}

// @Summary create a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param message body userCreateRequest true "user creation data"
// @Success 200 {object} APIError
// @Failure 400 {object} APIError
// @Failure 401 {object} APIError
// @Failure 404 {object} APIError
// @Failure 500 {object} APIError
// @Router /user [post]
// @Security ApiKeyAuth
func (u *UserHandler) create(w http.ResponseWriter, r *http.Request) {

	req := new(userCreateRequest)

	if err := render.Bind(r, req); err != nil {
		u.logger.WithError(err).Error("request body malformed")
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, "invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		_ = render.Render(w, r, newAPIError(http.StatusBadRequest, err.Error()))
		return
	}

	us := &billing.User{
		Tag:   req.Tag,
		Email: req.Email.String(),
	}

	if err := u.userRepo.Create(us); err != nil {
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	_ = render.Render(w, r, newAPIStatus(http.StatusOK, "successfully created a new user"))
}
