package msgservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/phamvinhdat/httpclient"
	"github.com/phamvinhdat/httpclient/body"
	"github.com/phamvinhdat/httpclient/hook"
	"github.com/sirupsen/logrus"
)

var sendMessageError = errors.New("can not send message")
var unexpectedStatusError = errors.New("unexpected http status code")

type (
	Service interface {
		Send(msg Message) error
	}

	service struct {
		msgServerURL string
		httpClient   httpclient.Client
	}
)

func New(msgServerURL string, httpClient httpclient.Client) Service {
	return &service{
		msgServerURL: msgServerURL,
		httpClient:   httpClient,
	}
}

func (s *service) Send(msg Message) error {
	broadcastURL := fmt.Sprintf("%s/broadcast/msg", s.msgServerURL)
	var res MsgServerRes
	statusCode, err := s.httpClient.Post(context.Background(), broadcastURL,
		httpclient.WithBodyProvider(body.NewJson(msg)),
		httpclient.WithHookFn(hook.UnmarshalResponse(&res)),
	)
	if err != nil {
		logrus.WithField("error", err).Error("failed to send message")
		return sendMessageError
	}
	if statusCode != http.StatusOK {
		logrus.Errorf("send message to server with unexpected http status code: %d, message: %s",
			statusCode, res.Message)
		return unexpectedStatusError
	}

	logrus.Info("send message to server success")
	return nil
}
