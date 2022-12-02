package config

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

// DatabaseConfiguration SQLite database configuration settings
type DatabaseConfiguration struct {
	Database string `mapstructure:"database"`
}

// Open attempts to open a new SQLite database connection based on the provided settings. Returns the connection and any errors encountered
func (me DatabaseConfiguration) Open() (*sql.DB, error) {

	return sql.Open("sqlite", me.Database)
}
