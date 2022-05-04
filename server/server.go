package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	applogger "github.com/fluidcoins/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/rs/cors"
	"github.com/yimikao/billing/core/oauth"
	"github.com/yimikao/billing/database/postgres"
	redisDB "github.com/yimikao/billing/database/redis"

	"golang.org/x/oauth2"
)

type ContextKey string

const (
	userCtx ContextKey = "user"
)

func New(
	port int64, cfg *oauth2.Config, logger applogger.Entry, dbConn *pg.DB) *http.Server {

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: SetupRoutes(dbConn, cfg, logger),
	}

}

func SetupRoutes(dbConn *pg.DB, cfg *oauth2.Config, logger applogger.Entry) http.Handler {

	router, oauthclient := chi.NewRouter(), oauth.NewGoogleOauthClient(cfg)

	router.Use(middleware.RequestID)
	//
	router.Use(middleware.AllowContentType("application/json"))
	//

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{ "status" : "ok" }`)
	})

	context := context.Background()

	// homepageHandler := NewHomepageHandler(nil, logger)
	// router.Get("/", homepageHandler)

	loginHandler := NewLoginHandler(cfg, logger)
	router.Get("/auth/google/login", loginHandler.Login)

	userLayer := postgres.NewUserLayer(dbConn)

	callbackHandler := NewCallbackHandler(oauthclient, userLayer, logger)
	router.Get("/auth/google/callback", callbackHandler.Callback)

	userRegistrationHandler := NewUserRegistrationHandler(userLayer, logger, &redisDB.Client{}, context)
	router.Post("/register", userRegistrationHandler.registerUser)

	return cors.AllowAll().Handler(router)
}

func Run(s *http.Server, entryLogger applogger.Entry) {

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
