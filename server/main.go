package main

import (
	"log"
	"os"

	"avocado.com/internal/dao"
	"avocado.com/internal/ginx"
	"avocado.com/internal/handler"
	"avocado.com/internal/model"
	"avocado.com/internal/service"
	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

func main() {
	c := parseConfig()

	d := dao.New(c.DatabaseConfig)
	s := service.New(c, d)
	handler.Init(s, c)

	r := gin.Default()
	ginx.Route(r, d, c.URLConfig)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func parseConfig() (c model.ServerConfig) {
	var path string
	if os.Getenv("ENV") == "production" {
		path = "config/production.yaml"
	} else {
		path = "config/development.yaml"
	}

	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := yaml.Unmarshal(f, &c); err != nil {
		log.Fatalf("Failed to parse config from file: %v", err)
	}
	return
}
