package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migrateCh "github.com/golang-migrate/migrate/v4/database/clickhouse"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed ch/*.sql
var migrationsChFS embed.FS

func runMigrationsClickhouse(connStr string) error {
	db, err := sql.Open("clickhouse", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	source, err := iofs.New(migrationsChFS, "ch")
	if err != nil {
		return err
	}

	driver, err := migrateCh.WithInstance(db, &migrateCh.Config{
		MigrationsTableEngine: "MergeTree",
	})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", source, "clickhouse", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func Run() error {
	clickhouseServer := os.Getenv("CLICKHOUSE_SERVER")
	clickhouseDatabase := os.Getenv("CLICKHOUSE_DATABASE")
	clickhouseUsername := os.Getenv("CLICKHOUSE_USERNAME")
	clickhousePassword := os.Getenv("CLICKHOUSE_PASSWORD")
	clickhouseTls := os.Getenv("CLICKHOUSE_TLS")

	tlsConfig := "&secure=true"
	if clickhouseTls == "false" {
		tlsConfig = ""
	}

	return runMigrationsClickhouse(fmt.Sprintf(`clickhouse://%s?username=%s&password=%s&database=%s%s`, clickhouseServer, clickhouseUsername, clickhousePassword, clickhouseDatabase, tlsConfig))
}
