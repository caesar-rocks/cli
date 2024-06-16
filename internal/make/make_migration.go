package make

import (
	"fmt"
	"os"
	"time"

	"github.com/caesar-rocks/cli/util"
)

func MakeMigration(migrationName string) error {
	migrationNameInSnakeCase := util.CamelToSnake(migrationName)

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
`, migrationName, timestamp, migrationName, timestamp, migrationName, timestamp, migrationName, timestamp)

	path := fmt.Sprintf("./database/migrations/%d_%s.go", timestamp, migrationNameInSnakeCase)

	err := os.WriteFile(path, []byte(migrationFileContents), 0644)
	if err != nil {
		return err
	}
	fmt.Println("Migration created successfully at", path)

	return nil
}
