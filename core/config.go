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
	consulKVPair   = "fluidcoins/flip"
)

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
