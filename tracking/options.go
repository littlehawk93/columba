package tracking

// Options defines configurable options for retrieving tracking information
type Options struct {
	UserAgent      string `mapstructure:"useragent"`
	TimeoutSeconds int    `mapstructure:"timeout"`
}
