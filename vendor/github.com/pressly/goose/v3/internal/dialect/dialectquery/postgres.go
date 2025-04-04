package dialectquery

import (
	"fmt"
)

type Postgres struct{}

var _ Querier = (*Postgres)(nil)

func (p *Postgres) CreateTable(tableName string) string {
	q := `CREATE TABLE %s (
		id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
		version_id bigint NOT NULL,
		is_applied boolean NOT NULL,
		tstamp timestamp NOT NULL DEFAULT now()
	)`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) InsertVersion(tableName string) string {
	q := `INSERT INTO %s (version_id, is_applied) VALUES ($1, $2)`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) DeleteVersion(tableName string) string {
	q := `DELETE FROM %s WHERE version_id=$1`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) GetMigrationByVersion(tableName string) string {
	q := `SELECT tstamp, is_applied FROM %s WHERE version_id=$1 ORDER BY tstamp DESC LIMIT 1`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) ListMigrations(tableName string) string {
	q := `SELECT version_id, is_applied from %s ORDER BY id DESC`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) GetLatestVersion(tableName string) string {
	q := `SELECT max(version_id) FROM %s`
	return fmt.Sprintf(q, tableName)
}

func (p *Postgres) TableExists(tableName string) string {
	schemaName, tableName := parseTableIdentifier(tableName)
	if schemaName != "" {
		q := `SELECT EXISTS ( SELECT 1 FROM pg_tables WHERE schemaname = '%s' AND tablename = '%s' )`
		return fmt.Sprintf(q, schemaName, tableName)
	}
	q := `SELECT EXISTS ( SELECT 1 FROM pg_tables WHERE (current_schema() IS NULL OR schemaname = current_schema()) AND tablename = '%s' )`
	return fmt.Sprintf(q, tableName)
}
