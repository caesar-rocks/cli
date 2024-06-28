package tools

import (
	"fmt"
	"os"
	"time"

	"github.com/caesar-rocks/cli/util"
)

type MakeMigrationOpts struct {
	MigrationName string `description:"The name of the migration to create"`
}

func (wrapper *ToolsWrapper) MakeMigration(opts MakeMigrationOpts) error {
	migrationNameInSnakeCase := util.CamelToSnake(opts.MigrationName)

	timestamp := time.Now().Unix()

	migrationFileContents := fmt.Sprintf(`package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func %sMigrationUp_%d(ctx context.Context, db *bun.DB) error {
	return nil
}

func %sMigrationDown_%d(ctx context.Context, db *bun.DB) error {
	return nil
}

func init() {
	Migrations.MustRegister(%sMigrationUp_%d, %sMigrationDown_%d)
}
`, opts.MigrationName, timestamp, opts.MigrationName, timestamp, opts.MigrationName, timestamp, opts.MigrationName, timestamp)

	path := fmt.Sprintf("./database/migrations/%d_%s.go", timestamp, migrationNameInSnakeCase)

	err := os.WriteFile(path, []byte(migrationFileContents), 0644)
	if err != nil {
		return err
	}
	fmt.Println("Migration created successfully at", path)

	return nil
}
