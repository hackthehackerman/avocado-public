package dao

import (
	"github.com/jmoiron/sqlx"
)

func (d *Dao) InsertNewOAuthStateToken(r *OauthStateToken, db *sqlx.DB) (err error) {
	if _, err = db.Exec(`INSERT INTO oauth_state_token VALUES(?,?,?,?,?)`, r.Id, r.Service, r.UserId, r.Token, r.CreatedAt); err != nil {
		return
	}
	return
}

func (d *Dao) GetOAuthStateToken(token, service string, db *sqlx.DB) (accessToken *OauthStateToken, err error) {
	var accessTokens []*OauthStateToken
	if err = db.Select(&accessTokens, "select * from oauth_state_token where token = ? and service = ? order by created_at desc", token, service); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}
