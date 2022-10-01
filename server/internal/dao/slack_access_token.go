package dao

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func (d *Dao) SaveSlackAccessToken(r *SlackAccessToken, db *sqlx.DB) (err error) {
	query := fmt.Sprintf(`INSERT INTO %s VALUES(?,?,?,?,?,?,?)`, slackAccessToken)
	if _, err = db.Exec(query, r.Id, r.UserId, r.TeamId, r.AccessToken, r.RefreshToken, r.ExpiredIn, r.CreatedAt); err != nil {
		return
	}
	return
}

func (d *Dao) GetSlackAccessToken(userId string, db *sqlx.DB) (accessToken *SlackAccessToken, err error) {
	var accessTokens []*SlackAccessToken
	if err = db.Select(&accessTokens, "select * from slack_access_token where user_id = ? order by created_at desc", userId); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}

func (d *Dao) GetSlackAccessTokenByTeamId(teamId string, db *sqlx.DB) (accessToken *SlackAccessToken, err error) {
	var accessTokens []*SlackAccessToken
	if err = db.Select(&accessTokens, "select * from slack_access_token where team_id = ? order by created_at desc", teamId); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}
