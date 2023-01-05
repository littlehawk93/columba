package tracking

import "database/sql"

// Migrate create any tables and data structure neccessary
func Migrate(db *sql.DB) error {

	if err := migrateV1(db); err != nil {
		return err
	}
	return nil
}

func migrateV1(db *sql.DB) error {

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS packages (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, status INTEGER NOT NULL, tracking_number TEXT NOT NULL, service TEXT NOT NULL, label TEXT NULL, created_on TEXT NOT NULL, last_updated_on TEXT NOT NULL, origin TEXT NULL, destination TEXT NULL, estimated_delivery_date TEXT NULL)"); err != nil {
		return err
	}

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS events (id TEXT PRIMARY KEY NOT NULL,package_id INTEGER NOT NULL,event_text TEXT NULL,details TEXT NULL,location TEXT NULL,event_timestamp TEXT NULL,FOREIGN KEY (package_id) REFERENCES packages(id));")
	return err
}
