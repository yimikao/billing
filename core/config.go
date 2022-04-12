package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	cfgSingleton   atomic.Value
	ErrKeyNotFound = fmt.Errorf("%s not found", consulKVPair)
	consulKVPair   = "yimikao/billing"
)

var (
	OauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://localhost:9090/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	OauthState = "state"
)

func InitOauth() {
	OauthConfig.ClientID = viper.GetString("CLIENT_ID")
	OauthConfig.ClientSecret = viper.GetString("CLIENT_SECRET")
	// oauthState = viper.GetString("oauthStateString")
}

func InitViper() {
	viper.SetConfigName("app")

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetConfigType("env")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

type Config struct {
	Database struct {
		Redis      string `json:"redis"`
		Postgresql string `json:"postgresql"`
	} `json:"database"`
}

func Load(client *api.Client) error {
	if _, ok := cfgSingleton.Load().(Config); ok {
		return nil
	}

	kv := client.KV()

	pair, _, err := kv.Get(consulKVPair, &api.QueryOptions{})
	if err != nil {
		return err
	}

	if pair == nil {
		return ErrKeyNotFound
	}

	var cfg Config

	if err := json.NewDecoder(bytes.NewBuffer(pair.Value)).Decode(&cfg); err != nil {
		return err
	}

	cfgSingleton.Store(cfg)
	return nil
}

func Set(c Config) {
	cfgSingleton.Store(c)
}

func Global() Config {
	cfg, ok := cfgSingleton.Load().(Config)
	if ok {
		return cfg
	}

	return Config{}
}
