package cmd

import (
	"context"
	"errors"
	"log"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
)

var deleteNoteCommand = &cobra.Command{
	Use:   "delete-note",
	Short: "Delete an existing note",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		noteID, err := cmd.Flags().GetInt64("id")
		if err != nil {
			log.Fatal(err)
		}

		if noteID == 0 {
			log.Fatal(errors.New("note ID required"))
		}

		notesRepo, err := notes.NewRepository(&notes.Config{DB: sqliteDB})
		if err != nil {
			log.Fatal(err)
		}

		err = notesRepo.Delete(ctx, noteID)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Note deleted.")
	},
}

func init() {
	rootCommand.AddCommand(deleteNoteCommand)

	deleteNoteCommand.Flags().Int64P("id", "i", 0, "The ID of the note to delete.")
}
