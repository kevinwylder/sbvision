package database

import (
	"database/sql"
	"fmt"
)

// SBDatabase is a namespace of database queries
type SBDatabase struct {
	*sql.DB
}

// ConnectToDatabase uses the DB_CREDS environment variable to connect to the database
func ConnectToDatabase(creds string) (*SBDatabase, error) {
	db, err := sql.Open("mysql", creds)
	if err != nil {
		return nil, fmt.Errorf("\n\tError connecting to the database: %s", err.Error())
	}
	return &SBDatabase{db}, nil

}
