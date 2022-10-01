package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (d *Dao) SaveHTTPRequest(r *HttpRequest, db *sqlx.DB) (err error) {
	query := fmt.Sprintf(`INSERT INTO %s VALUES(?,?,?,?,?,?)`, httpRequest)
	if _, err = db.Exec(query, r.Id, r.URI, r.Header, r.Body, r.OriginatedFromUs, r.CreatedAt); err != nil {
		return
	}
	return
}
