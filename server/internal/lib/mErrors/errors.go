package mErrors

import "net/http"

type Error struct {
	Code int
	Msg  string
}

func (e Error) Error() string {
	return e.Msg
}

var BadRequestError = Error{http.StatusBadRequest, "Bad Request Error"}
var InternalServerError = Error{http.StatusInternalServerError, "Internal Server Error"}
var UserNotFoundError = Error{http.StatusUnauthorized, "User not found"}
