package migrations

import "embed"

//go:embed postgres
var Postgres embed.FS

// Access the embedded migrations with:
// migrationsFS := &migrations.Postgres
