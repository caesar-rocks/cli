package main

import (
	"fmt"
	"os"
	"time"
	"unicode"

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
			tableName                string
		)
		huh.NewInput().Title("How is your migration named?").Value(&migrationName).Run()
		huh.NewInput().Title("What is the name of the table you want to make changes to?").Value(&tableName).Run()
		migrationNameInSnakeCase = camelToSnake(migrationName)

		timestamp := time.Now().Unix()

		migrationFileContents := fmt.Sprintf(`package migrations

import ormMigrations "caesar/orm/migrations"

func %sMigrationUp_%d(schema *ormMigrations.Schema) {
	schema.CreateTable("%s", func(t *ormMigrations.Table) {
		// Add your table columns here.
	})		
}

func %sMigrationDown_%d(schema *ormMigrations.Schema) {
	schema.DropTable("%s")
}

func init() {
	MigrationsStore.AddMigration(
		"%d_%s", 
		%sMigrationUp_%d, 
		%sMigrationDown_%d,
	)
}
`,
			migrationNameInSnakeCase,
			timestamp,
			tableName,
			migrationNameInSnakeCase,
			timestamp,
			tableName,
			timestamp,
			migrationNameInSnakeCase,
			migrationNameInSnakeCase,
			timestamp,
			migrationNameInSnakeCase,
			timestamp,
		)

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

func camelToSnake(s string) string {
	var result []rune

	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && unicode.IsLower(rune(s[i-1])) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
