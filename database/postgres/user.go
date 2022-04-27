package postgres

import (
	pg "github.com/go-pg/pg/v10"
	"github.com/yimikao/billing"
)

type userLayer struct {
	db *pg.DB
}

func NewUserLayer(db *pg.DB) billing.UserRepository {
	return &userLayer{db: db}
}

func (l *userLayer) Create(u *billing.User) error {

	ctx, cancelFn := WithContext()
	defer cancelFn()

	_, err := l.db.WithContext(ctx).Model(u).Insert()

	return err
}

func (l *userLayer) GetByReference(reference string) (*billing.User, error) {

	ctx, cancelFunc := WithContext()
	defer cancelFunc()

	user := new(billing.User)

	if err := l.db.WithContext(ctx).Model(user).
		Where("user.reference = ?", reference).
		Select(); err != nil {
		return user, err
	}

	return user, nil

}
