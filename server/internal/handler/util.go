package handler

import (
	"fmt"
	"net/http"
	"os"

	"avocado.com/internal/lib/mErrors"
	"github.com/gin-gonic/gin"
)

func returnWithJSON(c *gin.Context, obj interface{}, err error) {
	if err != nil {
		returnWithError(c, err)
	} else {
		c.JSON(200, obj)
	}
}

func returnWithError(c *gin.Context, err error) {
	if os.Getenv("ENV") != "production" {
		fmt.Println(fmt.Errorf("Error: %v", err))
	}

	if e, ok := err.(mErrors.Error); ok {
		c.JSON(e.Code, e.Error())
	} else {
		c.JSON(http.StatusInternalServerError, "")
	}
}
