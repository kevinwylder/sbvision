package database

import "database/sql"

// SBVisionDatabase is a namespace of database queries
type SBVisionDatabase struct {
	*sql.DB
}
