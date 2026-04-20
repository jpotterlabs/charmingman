package db

import (
	"database/sql"
	"fmt"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed" // Import to embed the SQLite WASM binary
	"github.com/ncruces/go-sqlite3/driver"
)

func openDB(dbPath string) (*sql.DB, error) {
	db, err := driver.Open(dbPath, func(c *sqlite3.Conn) error {
		for name, value := range pragmas {
			if err := c.Exec(fmt.Sprintf("PRAGMA %s = %s;", name, value)); err != nil {
				return fmt.Errorf("failed to set pragma %q: %w", name, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	return db, nil
}
