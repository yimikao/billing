package billing

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID `json:"id"`
	Reference string    `json:"reference"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func RunMigrations(db *pg.DB) error {
	models := []interface{}{

		(*User)(nil),
		(*Wallet)(nil),
		(*Transaction)(nil),
		(*TransactionEntry)(nil),
	}

	for _, m := range models {
		err := db.Model(m).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func RunSeeds(db *pg.DB) error {
	return nil
}
