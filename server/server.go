package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func New(port int64, h http.Handler) *http.Server {

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h,
	}

}

type ContextKey string

const (
	userCtx ContextKey = "user"
)

func SetupRoutes(h ...http.Handler) http.Handler {

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	//
	router.Use(middleware.AllowContentType("application/json"))
	//

	router.Get("health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "status" : "ok" }`)
	})

	return cors.AllowAll().Handler(router)
}
