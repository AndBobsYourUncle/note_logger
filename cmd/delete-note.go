package cmd

import (
	"context"
	"errors"
	"note-logger/internal/databases/sqlite"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
)

var deleteNoteCommand = &cobra.Command{
	Use:   "delete-note",
	Short: "Delete an existing note",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		noteID, err := cmd.Flags().GetInt64("id")
		if err != nil {
			cmd.PrintErr(err)
			return err
		}

		if noteID == 0 {
			err := errors.New("note ID required")
			cmd.PrintErr(err)
			return err
		}

		sqliteDB, err := sqlite.New(ctx)
		if err != nil {
			cmd.PrintErr(err)
			return err
		}

		notesRepo, err := notes.NewRepository(&notes.Config{DB: sqliteDB})
		if err != nil {
			cmd.PrintErr(err)
			return err
		}

		err = notesRepo.Delete(ctx, noteID)
		if err != nil {
			cmd.PrintErr(err)
			return err
		}

		cmd.Println("Note deleted.")

		return nil
	},
}

func init() {
	rootCommand.AddCommand(deleteNoteCommand)

	deleteNoteCommand.Flags().Int64P("id", "i", 0, "The ID of the note to delete.")
}
