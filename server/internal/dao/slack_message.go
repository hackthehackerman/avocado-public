package dao

import (
	"github.com/jmoiron/sqlx"
)

func (d *Dao) SaveSlackMessage(r *SlackMessage, db *sqlx.DB) (err error) {
	if _, err = db.Exec(`INSERT INTO slack_message VALUES(?,?,?,?,?,?,?,?,?,?)`, r.Id, r.SlackClientMessageId, r.SlackThreadTS, r.SlackTS, r.SlackChannelId, r.SlackTeamId, r.LinearIssueId, r.LinearCommentId, r.CreatedAt, r.Deleted); err != nil {
		return
	}
	return
}

func (d *Dao) UpdateSlackMessage(r *SlackMessage, db *sqlx.DB) (err error) {
	if _, err = db.Exec(`update slack_message set slack_client_message_id=?, slack_thread_ts=?, slack_ts=?, slack_channel_id=?, slack_team_id=?, linear_issue_id = ?, linear_comment_id = ?, created_at = ?, deleted = ? where id =?`,
		r.SlackClientMessageId, r.SlackThreadTS, r.SlackTS, r.SlackChannelId, r.SlackTeamId, r.LinearIssueId, r.LinearCommentId, r.CreatedAt, r.Deleted, r.Id); err != nil {
		return
	}
	return
}

func (d *Dao) GetSlackMessageByThreadTs(ts string, db *sqlx.DB) (s *SlackMessage, err error) {

	var accessTokens []*SlackMessage
	if err = db.Select(&accessTokens, "select * from slack_message where slack_thread_ts = ? order by created_at desc", ts); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}

func (d *Dao) GetSlackMessageByMessageTs(ts string, db *sqlx.DB) (s *SlackMessage, err error) {

	var accessTokens []*SlackMessage
	if err = db.Select(&accessTokens, "select * from slack_message where slack_ts = ? order by created_at desc", ts); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}

func (d *Dao) GetRootSlackMessageByLinearIssueId(issueId string, db *sqlx.DB) (s *SlackMessage, err error) {
	var accessTokens []*SlackMessage
	if err = db.Select(&accessTokens, "select * from slack_message where linear_issue_id = ? order by created_at ASC", issueId); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}

func (d *Dao) GetSlackMessageByLinearCommentId(commentId string, db *sqlx.DB) (s *SlackMessage, err error) {

	var accessTokens []*SlackMessage
	if err = db.Select(&accessTokens, "select * from slack_message where linear_comment_id = ? order by created_at desc", commentId); err != nil {
		return
	}
	if len(accessTokens) == 0 {
		return nil, err
	}
	return accessTokens[0], nil
}
