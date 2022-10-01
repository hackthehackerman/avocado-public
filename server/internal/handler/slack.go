package handler

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func HandleEventRequest(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	resp, err := s.ProcessSlackEvents(body, c.Request.Header)
	returnWithJSON(c, resp, err)
}

func HandleRedirectRequest(c *gin.Context) {
	code := c.Query("code")
	error := c.Query("error")
	states := c.Query("state")

	if error != "" {
		c.Status(http.StatusForbidden)
		return
	}

	_, err := s.ProcessSlackRedirect(code, states)
	if err != nil {
		returnWithError(c, err)
	}
	path := url.URL{Path: sc.URLConfig.Dashboard}
	c.Redirect(http.StatusFound, path.RequestURI())
}
