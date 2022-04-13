package cmd

import (
	"context"
	"log"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
)

var deleteNoteCommand = &cobra.Command{
	Use:   "delete-note",
	Short: "Delete an existing note",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		noteID, err := cmd.Flags().GetInt64("i")
		if err != nil {
			log.Fatal(err)
		}

		notesRepo, err := notes.NewRepository()
		if err != nil {
			log.Fatal(err)
		}

		err = notesRepo.Delete(ctx, noteID)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Note deleted.")

		err = notesRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCommand.AddCommand(deleteNoteCommand)

	deleteNoteCommand.PersistentFlags().Int64("i", 0, "The ID of the note to delete.")
}
