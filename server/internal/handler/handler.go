package handler

import (
	"avocado.com/internal/model"
	"avocado.com/internal/service"
)

var s *service.Service
var sc model.ServerConfig

func Init(service *service.Service, config model.ServerConfig) {
	s = service
	sc = config
}
