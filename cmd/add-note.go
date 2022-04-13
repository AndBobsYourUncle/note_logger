package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"note-logger/internal/databases/sqlite"
	"time"

	"note-logger/internal/entities"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
)

var addNoteCommand = &cobra.Command{
	Use:   "add-note",
	Short: "Add a new note",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		noteLine, err := cmd.Flags().GetString("content")
		if err != nil {
			log.Fatal(err)
		}

		if noteLine == "" {
			log.Fatal(errors.New("note content required"))
		}

		sqliteDB, err := sqlite.New(ctx)
		if err != nil {
			log.Fatal(err)
		}

		notesRepo, err := notes.NewRepository(&notes.Config{DB: sqliteDB})
		if err != nil {
			log.Fatal(err)
		}

		note, err := notesRepo.Create(ctx, &entities.Note{
			Content: noteLine,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Note added:\n%v - %v: %v\n", note.ID, note.CreatedAt.Format(time.Stamp), note.Content)
	},
}

func init() {
	rootCommand.AddCommand(addNoteCommand)

	addNoteCommand.Flags().StringP("content", "c", "", "The note contents to add.")
}
