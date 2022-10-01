package handler

import (
	"net/http"
	"net/url"

	"avocado.com/internal/lib/mErrors"
	"github.com/gin-gonic/gin"
)

func HandleGoogleRedirect(c *gin.Context) {
	csrf_cookie, err := c.Cookie("g_csrf_token")
	if csrf_cookie == "" || err != nil {
		returnWithError(c, mErrors.Error{Code: http.StatusBadRequest, Msg: "no csrf in cookie"})
		return
	}

	type Body struct {
		Credential string `form:"credential" binding:"required"`
		CSRF       string `form:"g_csrf_token" binding:"required"`
	}

	var body Body
	if err := c.Bind(&body); err != nil {
		returnWithError(c, mErrors.Error{Code: http.StatusBadRequest, Msg: "couldn't deserialize body"})
	}

	if body.CSRF != csrf_cookie {
		returnWithError(c, mErrors.Error{Code: http.StatusBadRequest, Msg: "failed to verify cookie"})
		return
	}

	token, err := s.ProcessGoogleRedirect(body.Credential)
	if err != nil {
		returnWithJSON(c, nil, err)
		return
	}

	c.SetCookie("session_token", token, 60*60*24*30, "/", "", true, true)
	path := url.URL{Path: sc.URLConfig.Dashboard}
	c.Redirect(http.StatusFound, path.RequestURI())
}

func HandleGetUserSettings(c *gin.Context) {
	userId := c.GetString("user_id")

	resp, err := s.GetUserSettings(userId)
	returnWithJSON(c, resp, err)
}
