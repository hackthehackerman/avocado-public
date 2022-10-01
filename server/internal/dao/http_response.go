package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (d *Dao) SaveHTTPResponse(r *HttpResponse, db *sqlx.DB) (err error) {
	query := fmt.Sprintf(`INSERT INTO %s VALUES(?,?,?,?,?,?,?)`, httpResponse)
	if _, err = db.Exec(query, r.Id, r.Header, r.Body, r.Status, r.OriginatedFromUs, r.RequestID, r.CreatedAt); err != nil {
		return
	}
	return
}
