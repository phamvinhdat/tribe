package try

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo_SomeCase(t *testing.T) {
	var retryCount uint
	tests := []struct {
		name               string
		expectedError      error
		expectedRetryCount uint
		fn                 doFn
	}{
		{
			name:               "first success",
			expectedRetryCount: 0,
			expectedError:      nil,
			fn: func() error {
				return nil
			},
		},
		{
			name:               "has error",
			expectedError:      assert.AnError,
			expectedRetryCount: 0,
			fn: func() error {
				return assert.AnError
			},
		},
		{
			name:               "limit retry success",
			expectedError:      nil,
			expectedRetryCount: 3,
			fn: func() error {
				if retryCount < 3 {
					return Continue
				}

				return nil
			},
		},
		{
			name:               "limit retry error",
			expectedError:      assert.AnError,
			expectedRetryCount: 3,
			fn: func() error {
				if retryCount >= 3 {
					return assert.AnError
				}

				return Continue
			},
		},
	}

	doer := New(
		WithTimeout(300*time.Millisecond),
		WithInterval(10*time.Millisecond),
		WithRetryFn(func(n uint, err error) {
			retryCount++
		}),
	)
	for _, test := range tests {
		retryCount = 0
		t.Run(test.name, func(t *testing.T) {
			actualErr := doer(test.fn)
			assert.Equal(t, test.expectedError, actualErr)
			assert.Equal(t, test.expectedRetryCount, retryCount)
		})
	}
}

func TestDo_Timeout(t *testing.T) {
	retryCount := 0
	doer := New(
		WithRetryFn(func(n uint, err error) {
			retryCount++
		}),
		WithTimeout(300*time.Millisecond),
		WithInterval(10*time.Millisecond))
	err := doer(func() error {
		return Continue
	})
	assert.Error(t, err)
	assert.Equal(t, err, TimeOutErr)
}

func TestDo_TimeoutWithRetry(t *testing.T) {
	retryCount := 0
	wantErr := assert.AnError
	retryFn := New(
		WithRetryFn(func(n uint, err error) {
			retryCount++
		}),
		WithRetryIf(func(e error) bool {
			if e == wantErr {
				return true
			}
			return false
		}),
		WithTimeout(300*time.Millisecond),
		WithInterval(10*time.Millisecond))
	err := retryFn(func() error {
		return wantErr
	})
	assert.Error(t, err)
	assert.Equal(t, TimeOutErr, err)
}
