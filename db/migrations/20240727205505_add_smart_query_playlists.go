package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddSmartQueryPlaylist, downAddSmartQueryPlaylist)
}

func upAddSmartQueryPlaylist(_ context.Context, tx *sql.Tx) error {
	_, err := tx.Exec(`
alter table playlist
	add column smart_query varchar null;
`)
	return err
}

func downAddSmartQueryPlaylist(_ context.Context, tx *sql.Tx) error {
	return nil
}
