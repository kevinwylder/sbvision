package mysqldb

import (
	// this package only supports mysql
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"reflect"
	"runtime"
)

// SBDatabase is a namespace of database queries
type SBDatabase struct {
	db *sql.DB

	// in search.go
	dataNearestRotation *sql.Stmt
}

// ConnectToDatabase waits for a sql connection then prepares queries for runtime
func ConnectToDatabase(creds string) (*SBDatabase, error) {
	db, err := sql.Open("mysql", creds)
	if err != nil {
		return nil, fmt.Errorf("\n\tError connecting to the database: %s", err.Error())
	}
	var counter int
	for err := db.Ping(); err != nil; counter++ {
		if counter > 30 {
			return nil, fmt.Errorf("Timed out pinging the database: %s", err.Error())
		}
	}

	sb := &SBDatabase{db: db}

	var prepFunctions = [](func() error){
		// in search.go
		sb.prepareDataNearestRotation,
	}

	for _, f := range prepFunctions {
		err = f()
		if err != nil {
			return nil, fmt.Errorf("\n\tError preparing query %s: %s", runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), err.Error())
		}
	}
	return sb, nil
}
