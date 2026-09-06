package blog

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// The migrations travel inside the binary, so the schema cannot fall out of
// step with the code that expects it.
//
//go:embed migrations/*.sql
var migrations embed.FS

// migrate applies every file in migrations/ that this database has not seen,
// in the order their names sort. The leading zeroes in those names are what
// makes that order right: sorting is character by character, so 0010 after
// 0009 only works while the width stays the same.
func migrate(db *sql.DB) error {
	_, err := db.Exec(`create table if not exists schema_migrations (
	    version    text not null primary key,
	    applied_at text not null default (datetime('now'))
	) strict`)
	if err != nil {
		return fmt.Errorf("журнал кестесі: %w", err)
	}

	names, err := fs.Glob(migrations, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		version := strings.TrimSuffix(path.Base(name), ".sql")

		var applied int
		if err := db.QueryRow(
			`select count(*) from schema_migrations where version = ?`, version).
			Scan(&applied); err != nil {
			return fmt.Errorf("%s: журналды оқу: %w", version, err)
		}
		if applied == 1 {
			continue
		}

		body, err := migrations.ReadFile(name)
		if err != nil {
			return err
		}

		// One transaction per file. A migration that breaks halfway leaves
		// nothing behind — neither half a schema nor a line in the journal —
		// so fixing the file and running again starts from a clean place.
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(body)); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: %w", version, err)
		}
		if _, err := tx.Exec(
			`insert into schema_migrations (version) values (?)`, version); err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: журналға жазу: %w", version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("%s: %w", version, err)
		}
	}
	return nil
}
