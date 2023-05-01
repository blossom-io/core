package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"core/db/migrations"

	sq "github.com/Masterminds/squirrel"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	schemaVersion = 1
	sourceName    = "iofs"
)

var driverName = "pgx"

type Postgres struct {
	DB *sql.DB
	Sq sq.StatementBuilderType
}

// New creates postgres.
func New(pgURL string) (*Postgres, error) {
	var (
		pg  Postgres
		err error
	)

	pg.Sq = sq.StatementBuilder.PlaceholderFormat(sq.Dollar) // pgx supports only dollar format $1, $2, etc

	pg.DB, err = sql.Open(driverName, pgURL)
	if err != nil {
		log.Fatalln("postgres - New:", err)
	}

	if err := pg.ensureSchema(pgURL); err != nil {
		return nil, err
	}

	result, err := pg.DB.Exec("SELECT now();")
	if err != nil {
		return nil, err
	}

	fmt.Println(result.RowsAffected())

	return &pg, nil
}

// ensureSchema runs migrations.
func (pg *Postgres) ensureSchema(pgURL string) error {
	migrationsFS := &migrations.Postgres

	src, err := iofs.New(migrationsFS, "postgres")
	if err != nil {
		return err
	}
	defer src.Close()

	driver, err := postgres.WithInstance(pg.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	// defer driver.Close()

	m, err := migrate.NewWithInstance(sourceName, src, "postgres", driver)
	if err != nil {
		return err
	}
	// defer m.Close()

	err = m.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (pg *Postgres) Close() error {
	return pg.DB.Close()
}
