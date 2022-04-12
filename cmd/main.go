package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yimikao/billing/core"
	"github.com/yimikao/billing/server"
)

func main() {
	// Initialize Viper across the application
	core.InitViper()

	// Initialize Logger across the application
	// logger.InitializeZapCustomLogger()

	// Initialize Oauth2 Services
	core.InitOauth()

	// Routes for the application
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("home"))
		return
	})
	http.HandleFunc("/login", server.NewLoginHandler().Login)
	http.HandleFunc("/callback", server.NewCallbackHandler().Callback)

	// logger.Log.Info("Started running on http://localhost:" + viper.GetString("port"))
	fmt.Println("serving...")
	log.Fatal(http.ListenAndServe(":9090", nil))

}
