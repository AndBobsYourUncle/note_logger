package cmd

import (
	"context"
	"errors"
	"note-logger/internal/databases/sqlite"
	"time"

	"note-logger/internal/entities"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
)

var addNoteCommand = &cobra.Command{
	Use:   "add-note",
	Short: "Add a new note",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		noteLine, err := cmd.Flags().GetString("content")
		if err != nil {
			return err
		}

		if noteLine == "" {
			err := errors.New("note content required")
			return err
		}

		sqliteDB, err := sqlite.New(ctx)
		if err != nil {
			return err
		}

		notesRepo, err := notes.NewRepository(&notes.Config{DB: sqliteDB})
		if err != nil {
			return err
		}

		note, err := notesRepo.Create(ctx, &entities.Note{
			Content: noteLine,
		})
		if err != nil {
			return err
		}

		cmd.Printf("Note added:\n%v - %v: %v\n", note.ID, note.CreatedAt.Format(time.Stamp), note.Content)

		return nil
	},
}

func init() {
	rootCommand.AddCommand(addNoteCommand)

	addNoteCommand.Flags().StringP("content", "c", "", "The note contents to add.")
}
