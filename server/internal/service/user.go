package service

import (
	"context"
	"time"

	"avocado.com/internal/dao"
	"avocado.com/internal/model"
	sessionManager "avocado.com/internal/session"

	"avocado.com/internal/lib/mErrors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/idtoken"
)

func (s *Service) ProcessGoogleRedirect(t string) (token string, err error) {
	var tx *sqlx.Tx
	tx, _ = s.dao.DB.Beginx()
	defer tx.Rollback()

	var payload *idtoken.Payload
	payload, err = idtoken.Validate(context.Background(), t, s.config.GoogleConfig.ClientID)
	if err != nil {
		return
	}

	type Claims struct {
		ID         string `mapstructure:"sub"`
		Email      string `mapstructure:"email"`
		FamilyName string `mapstructure:"family_name"`
		GivenName  string `mapstructure:"given_name"`
		PictureURL string `mapstructure:"picture"`
	}

	var claims Claims
	if err = mapstructure.Decode(payload.Claims, &claims); err != nil {
		return
	}

	var user *dao.User
	if user, err = s.dao.GetUserByEmail(claims.Email, tx); err != nil {
		return
	}

	ts := time.Now().Unix()
	if user != nil {
		// update user info if first time
		if user.FirstName == "" && user.LastName == "" {
			user.FirstName = claims.GivenName
			user.LastName = claims.FamilyName
			user.UpdatedAt = ts
		}
		if err = s.dao.UpdateUser(user, tx); err != nil {
			return
		}
	} else {
		user = &dao.User{
			Id:        uuid.NewString(),
			FirstName: claims.GivenName,
			LastName:  claims.FamilyName,
			Email:     claims.Email,
			CreatedAt: ts,
			UpdatedAt: ts,
		}
		if err = s.dao.InsertUser(user, tx); err != nil {
			return
		}
	}

	// create new token
	sessionManager := sessionManager.NewSessionManager(s.dao)
	if token, err = sessionManager.CreateNewSessionToken(user.Id); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

func (s *Service) GetUserSettings(userId string) (settings *model.UserSettingResponse, err error) {
	var user *dao.User
	if user, err = s.dao.GetUserById(userId, s.dao.DB); err != nil {
		return
	} else if user == nil {
		return nil, mErrors.UserNotFoundError
	}

	var slackStateToken string
	if slackStateToken, err = s.getStateToken(userId, "slack"); err != nil {
		return
	}
	var slackRedirectURI = "https://slack.com/oauth/v2/authorize" +
		"?scope=channels:history,chat:write,reactions:read,reactions:write,users:read&user_scope=" +
		"&redirect_uri=" +
		s.config.SlackConfig.RedirectURI +
		"&state=" + slackStateToken +
		"&client_id=" + s.config.SlackConfig.ClientID

	var slackAccessToken *dao.SlackAccessToken
	if slackAccessToken, err = s.dao.GetSlackAccessToken(userId, s.dao.DB); err != nil {
		return
	}
	slackConnected := slackAccessToken != nil

	var linearStateToken string
	if linearStateToken, err = s.getStateToken(userId, "linear"); err != nil {
		return
	}
	var linearRedirectURI = "https://linear.app/oauth/authorize" +
		"?client_id=" + s.config.Linearconfig.ClientID +
		"&redirect_uri=" + s.config.Linearconfig.RedirectURI +
		"&state=" + linearStateToken +
		"&response_type=code" +
		"&scope=read,write" +
		"&prompt=consent" +
		"&actor=application"

	var linearAccessToken *dao.LinearAccessToken
	if linearAccessToken, err = s.dao.GetLinearAccessToken(userId, s.dao.DB); err != nil {
		return
	}
	linearConnected := linearAccessToken != nil

	resp := model.UserSettingResponse{
		UserId:            userId,
		SlackRedirectURI:  slackRedirectURI,
		SlackConnected:    slackConnected,
		LinearRedirectURI: linearRedirectURI,
		LinearConnected:   linearConnected,
	}
	return &resp, err
}

func (s *Service) getStateToken(userId string, service string) (token string, err error) {
	var stateToken *dao.OauthStateToken
	if stateToken, err = s.dao.GetOAuthStateToken(userId, service, s.dao.DB); err != nil {
		return
	} else if stateToken != nil {
		ts := time.Now().Unix()
		if stateToken.CreatedAt+60*10 > ts {
			return stateToken.Token, nil
		}
	}

	// generate new token
	token = uuid.NewString()
	newToken := dao.OauthStateToken{
		Id:        uuid.NewString(),
		Service:   service,
		UserId:    userId,
		Token:     token,
		CreatedAt: time.Now().Unix(),
	}
	if s.dao.InsertNewOAuthStateToken(&newToken, s.dao.DB); err != nil {
		return
	}

	return token, nil
}

func (s *Service) getUserFromStateToken(service, token string) (user *dao.User, err error) {

	var state *dao.OauthStateToken
	if state, err = s.dao.GetOAuthStateToken(token, service, s.dao.DB); err != nil {
		return
	} else if state == nil {
		return
	}

	if user, err = s.dao.GetUserById(state.UserId, s.dao.DB); err != nil {
		return
	} else if user == nil {
		return
	}

	return
}
