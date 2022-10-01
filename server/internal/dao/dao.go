package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"avocado.com/internal/model"
)

// wrapper generic for sqlx.DB & sqlx.Tx
type DBX interface {
	Exec(query string, args ...any) (sql.Result, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

type Dao struct {
	DB *sqlx.DB
}

func New(c model.DatabaseConfig) *Dao {
	db := sqlx.MustConnect("mysql", c.URI)

	d := &Dao{
		DB: db,
	}

	return d
}
