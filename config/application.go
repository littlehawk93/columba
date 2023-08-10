package config

import "github.com/littlehawk93/columba/tracking"

// ApplicationConfiguration full application configuration settings
type ApplicationConfiguration struct {
	Database                    *DatabaseConfiguration `mapstructure:"database"`
	Web                         *WebConfiguration      `mapstructure:"web"`
	Browser                     *tracking.Options      `mapstructure:"browser"`
	WebRoot                     string                 `mapstructure:"webroot"`
	MinimumRefreshTimeSeconds   int                    `mapstructure:"minrefreshtime"`
	BackgroundUpdateTimeSeconds int                    `mapstructure:"bgupdatetime"`
}
