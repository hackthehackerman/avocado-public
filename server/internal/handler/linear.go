package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HandleLinearWebhook(c *gin.Context) {
	body, _ := io.ReadAll(c.Request.Body)

	resp, err := s.ProcessLinearWebhook(body)
	returnWithJSON(c, resp, err)
}

func HandleLinearRedirectRequest(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.Status(http.StatusForbidden)
		return
	}

	_, err := s.ProcessLinearRedirect(code, state)
	if err != nil {
		returnWithError(c, err)
	}
	path := url.URL{Path: sc.URLConfig.Dashboard}
	c.Redirect(http.StatusFound, path.RequestURI())
}
