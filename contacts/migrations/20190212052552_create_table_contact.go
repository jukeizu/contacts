package migration

import (
	"database/sql"
)

type CreateTableContact20190212052552 struct{}

func (m CreateTableContact20190212052552) Version() string {
	return "20190212052552_CreateTableContact"
}

func (m CreateTableContact20190212052552) Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS contact (
			id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
			name STRING NOT NULL DEFAULT '',
			address STRING NOT NULL DEFAULT '',
			phone STRING NOT NULL DEFAULT '',
			created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated TIMESTAMPTZ
		)`)
	return err
}

func (m CreateTableContact20190212052552) Down(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE contact`)
	return err
}
