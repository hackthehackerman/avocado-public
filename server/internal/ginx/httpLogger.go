package ginx

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"avocado.com/internal/dao"
)

type HttpLogger struct {
	dao *dao.Dao
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func NewHttpLogger(d *dao.Dao) HttpLogger {
	logger := HttpLogger{
		dao: d,
	}
	return logger
}

func (h *HttpLogger) Log(c *gin.Context) {
	requestId := uuid.NewString()
	ts := time.Now().Unix()

	uri := c.Request.URL.String()

	header, _ := json.Marshal(c.Request.Header)
	headerString := string(header)

	body, _ := io.ReadAll(c.Request.Body)
	bodyString := string(body)
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	request := &dao.HttpRequest{
		Id:               requestId,
		URI:              &uri,
		Header:           &headerString,
		Body:             &bodyString,
		OriginatedFromUs: false,
		CreatedAt:        ts,
	}

	h.dao.SaveHTTPRequest(request, h.dao.DB)

	bw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bw

	// before request
	c.Next()
	// after request

	responseId := uuid.NewString()
	ts = time.Now().Unix()

	header, _ = json.Marshal(c.Writer.Header())
	headerString = string(header)

	bodyString = bw.body.String()

	status := c.Writer.Status()
	response := &dao.HttpResponse{
		Id:               responseId,
		Header:           &headerString,
		Body:             &bodyString,
		Status:           &status,
		OriginatedFromUs: true,
		RequestID:        &requestId,
		CreatedAt:        ts,
	}

	h.dao.SaveHTTPResponse(response, h.dao.DB)
}
