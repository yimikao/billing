package datastore

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"
)

func New(dsn string) (*pg.DB, error) {

	opts, err := pg.ParseURL(dsn)

	if err != nil {
		return nil, err
	}

	db := pg.Connect(opts)

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return db, nil

}

func withContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*4)
}
