package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync/atomic"

	"github.com/hashicorp/consul/api"
)

var (
	cfgSingleton   atomic.Value
	ErrKeyNotFound = fmt.Errorf("%s not found", consulKVPair)
	consulKVPair   = "yimikao/billing"
)

type Config struct {
	Oauth struct {
		OauthClientID     string   `json:"oauth_client_id"`
		OauthClientSecret string   `json:"oauth_client_secret"`
		OauthCallbackURL  string   `json:"oauth_callback_url"`
		OauthScopes       []string `json:"oauth_scopes"`
	} `json:"oauth"`

	Database struct {
		Postgresql string `json:"postgresql"`
		Redis      struct {
			Password string `json:"password"`
			Addr     string `json:"addr"`
			UseTLS   bool   `json:"use_tls"`
		} `json:"redis"`
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
