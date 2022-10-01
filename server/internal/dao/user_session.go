package dao

import (
	"time"
)

func (d *Dao) GetActiveSessionByUserId(uid string, db DBX) (userSession *UserSession, err error) {
	var userSessions []UserSession

	if err = db.Select(&userSessions, "select * from user_session where user_id = ? and deleted = false order by created_at desc", uid); err != nil {
		return
	}
	if len(userSessions) == 0 {
		return nil, err
	}

	return &userSessions[0], err
}

func (d *Dao) GetSessionByToken(token string, db DBX) (userSession *UserSession, err error) {
	var userSessions []UserSession

	if err = db.Select(&userSessions, "select * from user_session where token = ? order by created_at desc", token); err != nil {
		return
	}
	if len(userSessions) == 0 {
		return nil, err
	}

	return &userSessions[0], err
}

func (d *Dao) InvalidateUserSession(userSession *UserSession, db DBX) (err error) {
	updatedAt := time.Now().Unix()
	if _, err := db.Exec("update user_session set updated_at=?, deleted=true where id=?",
		updatedAt, userSession.Id); err != nil {

		return err
	}
	return
}

func (d *Dao) InsertNewUserSession(userSession *UserSession, db DBX) (err error) {
	if _, err = db.Exec("INSERT INTO user_session VALUES(?,?,?,?,?,?,?)", userSession.Id, userSession.UserId, userSession.Token, userSession.ExpiredAt, userSession.CreatedAt, userSession.UpdatedAt, userSession.Deleted); err != nil {
		return
	}
	return
}
