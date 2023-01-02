package config

// WebConfiguration web server configuration settings
type WebConfiguration struct {
	BindAddress string `mapstructure:"bind"`
	Port        uint16 `mapstructure:"port"`
}
