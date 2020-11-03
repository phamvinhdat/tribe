package try

import "time"

type (
	// Function signature of retry if function
	RetryIfFn func(error) bool

	// Function signature of WithRetryFn function
	// n = count of attempts
	OnRetryFn func(n uint, err error)

	opt struct {
		timeout, interval time.Duration
		onRetryFn         OnRetryFn
		retryIfFn         RetryIfFn
	}

	Option interface {
		apply(*opt)
	}

	optionFn func(*opt)
)

func (f optionFn) apply(args *opt) {
	f(args)
}

func WithTimeout(timeout time.Duration) Option {
	return optionFn(func(args *opt) {
		args.timeout = timeout
	})
}

func WithInterval(interval time.Duration) Option {
	return optionFn(func(args *opt) {
		args.interval = interval
	})
}

func WithRetryFn(onRetryFn OnRetryFn) Option {
	return optionFn(func(args *opt) {
		args.onRetryFn = onRetryFn
	})
}

func WithRetryIf(retryIfFn RetryIfFn) Option {
	return optionFn(func(args *opt) {
		args.retryIfFn = retryIfFn
	})
}
