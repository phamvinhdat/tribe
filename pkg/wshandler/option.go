package wshandler

type (
	Option interface {
		apply(*option)
	}

	option struct {
		onReMsgFn OnReceiveMessageFn
		onDisFn   OnDisconnectFn
	}

	optionFn func(*option)
)

func (optFn optionFn) apply(opt *option) {
	optFn(opt)
}

func OnReceiveMessage(onReMsgFn OnReceiveMessageFn) Option {
	return optionFn(func(opt *option) {
		opt.onReMsgFn = onReMsgFn
	})
}

func OnDisConnect(onDisFn OnDisconnectFn) Option {
	return optionFn(func(opt *option) {
		opt.onDisFn = onDisFn
	})
}
