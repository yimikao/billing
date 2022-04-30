package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	applogger "github.com/fluidcoins/log"
	"github.com/hashicorp/consul/api"
	"github.com/yimikao/billing/core"
	"github.com/yimikao/billing/core/oauth"
	"github.com/yimikao/billing/server"
)

func main() {

	var consulAddr string

	flag.StringVar(&consulAddr, "consul", "localhost:8500", "Url to the running consul instance")

	flag.Parse()

	cfg := api.DefaultConfig()

	cfg.Address = consulAddr

	fmt.Println(consulAddr)

	consulClient, err := api.NewClient(cfg)
	if err != nil {
		log.Fatalf("could not init consul client... %v", err)
	}

	if err := core.Load(consulClient); err != nil {
		log.Fatalf("Could not init configuration... %v", err)
	}

	conf := core.Global()

	// Initialize Logger across the application
	logger := applogger.New(applogger.LevelDebug, 4)
	// logger.InitializeZapCustomLogger()

	hostName, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	entryLogger := logger.WithFields(map[string]interface{}{
		"host": hostName,
		"app":  "billing",
	})

	oauthConfig := oauth.NewGoogleOauthConfig(conf)

	svr := server.New(8080, oauthConfig, entryLogger, nil)

	server.Run(svr, entryLogger)

}
