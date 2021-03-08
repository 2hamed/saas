package trace

type tracingConfig struct {
	traceName string
}

type TraceOpts func(*tracingConfig)

func WithName(name string) TraceOpts {
	return func(c *tracingConfig) {
		c.traceName = name
	}
}
