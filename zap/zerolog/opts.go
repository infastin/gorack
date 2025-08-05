package zapzerolog

type Option func(*core)

func WithCopyTimestamp(value bool) Option {
	return func(c *core) {
		c.copyTimestamp = value
	}
}

func WithCopyCaller(value bool) Option {
	return func(c *core) {
		c.copyCaller = value
	}
}

func WithCopyStack(value bool) Option {
	return func(c *core) {
		c.copyStack = value
	}
}
