package logger

type Option func(o *logger)

func WithTracing() Option {
	return func(o *logger) {
		o.Tracing = true
	}
}
