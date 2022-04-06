package datastore

import (
	"github.com/go-pg/pg/v10"
	"github.com/yimikao/billing"
)

type userLayer struct {
	db *pg.DB
}

func NewUserLayer(db *pg.DB) billing.UserRepository {
	return &userLayer{db: db}
}

func (l *userLayer) Create(u *billing.User) error {

	ctx, cancelFn := withContext()
	defer cancelFn()

	_, err := l.db.WithContext(ctx).Model(u).Insert()

	return err
}
