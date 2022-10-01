package service

import (
	"avocado.com/internal/dao"
	"avocado.com/internal/model"
)

type Service struct {
	config model.ServerConfig
	// slack  *slack.Client
	dao *dao.Dao
	c   chan func()
}

func New(c model.ServerConfig, dao *dao.Dao) *Service {
	// slackClient := slack.New("YOUR_TOKEN_HERE")

	s := &Service{
		config: c,
		// slack:  slackClient,
		dao: dao,
		c:   make(chan func()),
	}

	go worker(s.c)

	return s
}

func worker(c chan func()) {
	for {
		f := <-c
		go f()
	}
}
