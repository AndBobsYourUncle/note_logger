package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

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

func main() {
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

	if len(os.Args) < 3 {
		log.Fatal("Not enough arguments given")
	}

	commandName, err := getCommand()
	if err != nil {
		log.Fatal(err)
	}

	switch commandName {
	case addNote:
		noteLine := strings.Join(os.Args[2:], " ")

		_, insertErr := db.Exec(insertNote, time.Now(), noteLine)
		if insertErr != nil {
			log.Fatal(insertErr)
		}

		log.Println("Log added")
		log.Println(noteLine)
	case listRecent:
	}
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

func getCommand() (command, error) {
	commandName := command(os.Args[1])

	if commandName != addNote && commandName != listRecent {
		return "", errors.New("invalid command name")
	}

	return commandName, nil
}
