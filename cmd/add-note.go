package cmd

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "notes.sqlite"

const createTable string = `
  CREATE TABLE IF NOT EXISTS notes (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  note TEXT
  );`

const createIndex string = `
CREATE INDEX IF NOT EXISTS time_index 
ON notes(time);
`

const insertNote string = `
INSERT INTO notes VALUES(NULL,?,?);
`

type command string

const (
	addNote    command = "add-note"
	listRecent command = "list-recent"
)

var addNoteCommand = &cobra.Command{
	Use:   "add-note",
	Short: "Add a new note",
	Run: func(cmd *cobra.Command, args []string) {
		noteLine, err := cmd.Flags().GetString("n")
		if err != nil {
			log.Fatal(err)
		}

		ex, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		exPath := filepath.Dir(ex)

		filename := exPath + "/" + dbFile

		touchDBFile(filename)

		db, err := sql.Open("sqlite3", filename)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = db.Exec(createTable); err != nil {
			log.Fatal(err)
		}

		if _, err = db.Exec(createIndex); err != nil {
			log.Fatal(err)
		}

		_, insertErr := db.Exec(insertNote, time.Now(), noteLine)
		if insertErr != nil {
			log.Fatal(insertErr)
		}

		log.Println("Log added")
		log.Println(noteLine)
	},
}

func init() {
	rootCommand.AddCommand(addNoteCommand)

	addNoteCommand.PersistentFlags().String("n", "", "The note contents to add.")
}

func touchDBFile(filename string) {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		file, createErr := os.Create(filename)
		if createErr != nil {
			log.Fatal(err)
		}

		closeErr := file.Close()
		if closeErr != nil {
			log.Fatal(closeErr)
		}
	}
}
