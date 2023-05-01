package postgres

import (
	"testing"

	"core/db/migrations"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	st "github.com/golang-migrate/migrate/v4/source/testing"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Test(t *testing.T) {
	// reuse the embed.FS set in example_test.go
	migrationsFS := &migrations.Postgres

	d, err := iofs.New(migrationsFS, "postgres")
	if err != nil {
		t.Fatal(err)
	}

	st.Test(t, d)
}
