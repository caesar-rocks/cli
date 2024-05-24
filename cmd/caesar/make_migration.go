package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caesar-rocks/cli/util"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var makeMigrationCmd = &cobra.Command{
	Use:     "make:migration",
	Short:   "Create a new migration",
	GroupID: "make",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			migrationName            string
			migrationNameInSnakeCase string
		)

		if len(args) > 0 {
			migrationName = args[0]
		} else {
			huh.NewInput().Title("How is your migration named?").Value(&migrationName).Run()
		}

		migrationNameInSnakeCase = util.CamelToSnake(migrationName)

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
			panic(err)
		}
		fmt.Println("Migration created successfully at", path)
	},
}

func init() {
	rootCmd.AddCommand(makeMigrationCmd)
}
