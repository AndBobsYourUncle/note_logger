package cmd

import (
	"database/sql"
	"log"
	"note-logger/internal/databases/sqlite"

	"github.com/spf13/cobra"
)

// used by any command to initialize repositories
var sqliteDB *sql.DB

var rootCommand = &cobra.Command{
	Use: "note-logger",
}

func Execute() {
	db, err := sqlite.New()
	if err != nil {
		log.Fatal(err)
	}

	sqliteDB = db

	if err = rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
