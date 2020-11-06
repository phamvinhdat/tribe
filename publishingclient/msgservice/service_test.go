package msgservice

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	clientmock "github.com/phamvinhdat/httpclient/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	httpClient := clientmock.NewMockClient(ctrl)
	msgServerURL := "message server url"
	url := fmt.Sprintf("%s/broadcast/msg", msgServerURL)

	sendMessage := Message{
		Message:   "this is message",
		Timestamp: time.Now(),
	}

	tests := []struct {
		name          string
		prepare       func()
		expectedError error
	}{
		{
			name: "send message error",
			prepare: func() {
				httpClient.EXPECT().
					Post(context.Background(), url, gomock.Any()).
					Return(0, errors.New("return error"))
			},
			expectedError: sendMessageError,
		},
		{
			name: "unexpected status error",
			prepare: func() {
				httpClient.EXPECT().
					Post(context.Background(), url, gomock.Any()).
					Return(http.StatusNotFound, nil)
			},
			expectedError: unexpectedStatusError,
		},
		{
			name: "success",
			prepare: func() {
				httpClient.EXPECT().
					Post(context.Background(), url, gomock.Any()).
					Return(http.StatusOK, nil)
			},
			expectedError: nil,
		},
	}

	service := New(msgServerURL, httpClient)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.prepare()
			actualErr := service.Send(sendMessage)
			assert.Equal(t, test.expectedError, actualErr)
		})
	}
}
