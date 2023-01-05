package config

// ApplicationConfiguration full application configuration settings
type ApplicationConfiguration struct {
	Database                    *DatabaseConfiguration `mapstructure:"database"`
	Web                         *WebConfiguration      `mapstructure:"web"`
	WebRoot                     string                 `mapstructure:"webroot"`
	MinimumRefreshTimeSeconds   int                    `mapstructure:"minrefreshtime"`
	BackgroundUpdateTimeSeconds int                    `mapstructure:"bgupdatetime"`
}
