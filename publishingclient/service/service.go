package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/phamvinhdat/tribe/pkg/try"
	"github.com/phamvinhdat/tribe/publishingclient/msgservice"
)

type service struct {
	tryer      try.Doer
	msgService msgservice.Service
}

func New(tryer try.Doer, msgService msgservice.Service) *service {
	return &service{
		tryer:      tryer,
		msgService: msgService,
	}
}

func (s *service) Run() {
	_ = s.tryer(func() error {
		// random string
		str := uuid.New().String()
		if err := s.msgService.Send(msgservice.Message{
			Message:   str,
			Timestamp: time.Now(),
		}); err != nil {
			// handle error
		}

		return try.Continue // interval
	})
}
