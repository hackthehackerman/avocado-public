package ginx

import (
	"net/http"

	"avocado.com/internal/dao"
	"avocado.com/internal/model"
	sessionManager "avocado.com/internal/session"
	"github.com/gin-gonic/gin"
)

type SessionAuth struct {
	config         model.URLConfig
	sessionManager *sessionManager.SessionManager
}

func NewSessionAuth(d *dao.Dao, c model.URLConfig) SessionAuth {
	sm := sessionManager.NewSessionManager(d)
	a := SessionAuth{
		config:         c,
		sessionManager: &sm,
	}
	return a
}

func (a *SessionAuth) Auth(c *gin.Context) {
	var abortAndRedirect = func() {
		c.JSON(http.StatusUnauthorized, "you must be logged in")
		c.Abort()
	}

	t, err := c.Cookie("session_token")
	if t == "" || err != nil {
		abortAndRedirect()
		return
	}

	var userSession *dao.UserSession
	if userSession, err = a.sessionManager.GetValidSession(t); err != nil {
		abortAndRedirect()
		return
	} else if userSession == nil {
		abortAndRedirect()
		return
	}

	c.Set("user_id", userSession.UserId)
	c.Next()
}
