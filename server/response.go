package server

import (
	"net/http"

	"github.com/go-chi/render"
)

type meta struct {
	Paging pagingInfo `json:"paging"`
}

type pagingInfo struct {
	Total   int64 `json:"total"`
	PerPage int64 `json:"per_page"`
	Page    int64 `json:"page"`
}

type APIStatus struct {
	statusCode int
	Status     bool `json:"status"`
	// Generic message that tells you the status of the operation
	Message string `json:"message"`
}

func (a APIStatus) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, a.statusCode)
	return nil
}

type APIError struct {
	APIStatus
}

func newAPIStatus(code int, s string) APIStatus {
	return APIStatus{
		statusCode: code,
		Status:     true,
		Message:    s,
	}
}

func newAPIError(code int, s string) APIError {
	return APIError{
		APIStatus: APIStatus{
			statusCode: code,
			Status:     false,
			Message:    s,
		},
	}
}

var (
	errInvalidRequestBody = newAPIError(http.StatusBadRequest, "Invalid request body")
	errUnauthorized       = newAPIError(http.StatusUnauthorized, "You are not authorized to make this request")
	errForbidden          = newAPIError(http.StatusForbidden, "You are forbidden from making this request")
)

type UserResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}
