package store

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disabled")
	if err != nil {
		return nil, fmt.Errorf("db: open %w\n", err)
	}
	fmt.Println("Connected to db...")
	return db, nil
}
