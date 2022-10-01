package dao

import (
	"github.com/jmoiron/sqlx"
)

func (d *Dao) SaveLinearAccessToken(r *LinearAccessToken, db *sqlx.DB) (err error) {
	if _, err = db.Exec(`INSERT INTO linear_access_token VALUES(?,?,?,?,?,?)`, r.Id, r.UserId, r.AccessToken, r.RefreshToken, r.ExpiredAt, r.CreatedAt); err != nil {
		return
	}
	return
}

func (d *Dao) GetLinearAccessToken(userId string, db *sqlx.DB) (accessToken *LinearAccessToken, err error) {
	var accessTokens []*LinearAccessToken
	if err = db.Select(&accessTokens, "select * from linear_access_token where user_id = ? order by created_at desc", userId); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}
