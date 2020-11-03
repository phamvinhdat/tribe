package try

import (
	"errors"
	"time"
)

var (
	Continue   = errors.New("continue without error")
	TimeOutErr = errors.New("timeout")
)

type (
	// Doer will automatically retry doFn until successful or time out.
	Doer func(doFn doFn) (err error)

	doFn func() error
)

const (
	defaultTimeout  = 30 * time.Second
	defaultInterval = 1 * time.Second
)

// New return a Doer
func New(opts ...Option) Doer {
	opt := opt{
		timeout:   defaultTimeout,
		interval:  defaultInterval,
		onRetryFn: func(uint, error) {},
		retryIfFn: func(error) bool { return false },
	}
	for _, o := range opts {
		o.apply(&opt)
	}

	return func(fn doFn) (err error) {
		timer := time.NewTimer(opt.timeout)
		ticker := time.NewTicker(opt.interval)
		defer timer.Stop()
		defer ticker.Stop()
		var count uint
		for {
			err = fn()
			if err == nil {
				return
			}

			if err != Continue {
				if !opt.retryIfFn(err) {
					return
				}
			}

			select {
			case <-timer.C:
				err = TimeOutErr
				return
			case <-ticker.C:
				count++
				// on retry hook
				opt.onRetryFn(count, err)
				continue
			}
		}
	}
}
