package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fluidcoins/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/yimikao/billing/core/oauth"
	"golang.org/x/oauth2"
)

type ContextKey string

const (
	userCtx ContextKey = "user"
)

func New(port int64, h http.Handler) *http.Server {

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h,
	}

}

func SetupRoutes(cfg *oauth2.Config, logger log.Entry) http.Handler {

	router, oauthclient := chi.NewRouter(), oauth.NewGoogleOauthClient(cfg)

	router.Use(middleware.RequestID)
	//
	router.Use(middleware.AllowContentType("application/json"))
	//

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "status" : "ok" }`)
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "status" : "welcome home" }`)
	})

	loginHandler := NewLoginHandler(cfg, logger)
	router.Get("/auth/google/login", loginHandler.Login)

	callbackHandler := NewCallbackHandler(oauthclient, logger)
	router.Get("/auth/google/callback", callbackHandler.Callback)

	return cors.AllowAll().Handler(router)
}

func Run(s *http.Server, entryLogger log.Entry) {

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		entryLogger.Info("Starting HTTP server")
		if err := s.ListenAndServe(); err != nil {
			entryLogger.WithError(err).Error("ListenAndServe")
		}
	}()

	<-sig
	entryLogger.Debug("Shutting down server")
	if err := s.Shutdown(context.Background()); err != nil {
		entryLogger.WithError(err).Error("Could not shut down server properly")
	}

}
