package zapzerolog

type config struct {
	nameKey       string
	copyTimestamp bool
	copyCaller    bool
	copyStack     bool
}

func defaultConfig() *config {
	return &config{
		nameKey:       "",
		copyTimestamp: false,
		copyCaller:    false,
		copyStack:     false,
	}
}

type Option func(*config)

func WithNameKey(value string) Option {
	return func(c *config) {
		c.nameKey = value
	}
}

func WithCopyTimestamp(value bool) Option {
	return func(c *config) {
		c.copyTimestamp = value
	}
}

func WithCopyCaller(value bool) Option {
	return func(c *config) {
		c.copyCaller = value
	}
}

func WithCopyStack(value bool) Option {
	return func(c *config) {
		c.copyStack = value
	}
}
