package migrator

import (
	"database/sql"
	"fmt"
	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
	"os"
)

type Migrator struct {
	migrations   *migrate.FileMigrationSource
	dbConnection *sql.DB
	dialect      string
}

func NewMigrator(dbConnection *sql.DB, dialect string) Migrator {

	// create new connection to database
	migrations := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{migrations: migrations, dbConnection: dbConnection, dialect: dialect}
}

func (m Migrator) Up() {

	n, err := migrate.Exec(m.dbConnection, m.dialect, m.migrations, migrate.Up)
	if err != nil {

		panic(fmt.Sprintf("Migration error: %v", err))
	}

	fmt.Printf("Applied %d migrations!\n", n)
}

func (m Migrator) Down() {

	n, err := migrate.Exec(m.dbConnection, m.dialect, m.migrations, migrate.Down)
	if err != nil {

		panic(fmt.Sprintf("Migration error: %v", err))
	}

	fmt.Printf("Rollbacked %d migrations!\n", n)
}

func (m Migrator) Status() {

	fmt.Print("migration status: \n")

	records, err := migrate.GetMigrationRecords(m.dbConnection, m.dialect)
	if err != nil {
		panic(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "AppliedAt"})

	for _, val := range records {
		table.Append([]string{
			val.Id,
			val.AppliedAt.String(),
		})
	}

	table.Render()
}
