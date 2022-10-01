package sessionManager

import (
	"time"

	"avocado.com/internal/dao"
	"github.com/google/uuid"
)

type SessionManager struct {
	dao *dao.Dao
}

func NewSessionManager(d *dao.Dao) SessionManager {
	s := SessionManager{
		dao: d,
	}
	return s
}

func (s *SessionManager) GetActiveSessionToken(userId string) (token string, err error) {
	userSession, err := s.dao.GetActiveSessionByUserId(userId, s.dao.DB)
	if err != nil {
		return "", err
	}
	ts := time.Now().Unix()
	if userSession.ExpiredAt <= ts {
		return "", nil
	}

	return userSession.Token, nil
}

func (s *SessionManager) CreateNewSessionToken(userId string) (token string, err error) {
	tx, _ := s.dao.DB.Beginx()
	defer tx.Rollback()

	// invalidate existing session if needed
	var userSession *dao.UserSession
	userSession, err = s.dao.GetActiveSessionByUserId(userId, tx)
	if err != nil {
		return "", err
	}
	if userSession != nil {
		err = s.dao.InvalidateUserSession(userSession, tx)
		if err != nil {
			return "", err
		}
	}

	// new token
	ts := time.Now().Unix()
	id := uuid.NewString()
	token = uuid.NewString()
	newUserSession := dao.UserSession{
		Id:        id,
		UserId:    userId,
		Token:     token,
		ExpiredAt: ts + 60*60*24*30, // 30 days
		CreatedAt: ts,
		UpdatedAt: ts,
		Deleted:   false,
	}
	if err = s.dao.InsertNewUserSession(&newUserSession, tx); err != nil {
		return
	}
	tx.Commit()

	return newUserSession.Token, nil
}

func (s *SessionManager) GetValidSession(token string) (userSession *dao.UserSession, err error) {
	if userSession, err = s.dao.GetSessionByToken(token, s.dao.DB); err != nil {
		return
	}

	ts := time.Now().Unix()
	if userSession.ExpiredAt <= ts {
		return nil, nil
	}

	if userSession.Deleted {
		return nil, nil
	}
	return
}
