package cmd

import (
	"context"
	"log"
	"note-logger/internal/entities"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"

	_ "github.com/mattn/go-sqlite3"
)

var addNoteCommand = &cobra.Command{
	Use:   "add-note",
	Short: "Add a new note",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		noteLine, err := cmd.Flags().GetString("n")
		if err != nil {
			log.Fatal(err)
		}

		notesRepo, err := notes.NewRepository()
		if err != nil {
			log.Fatal(err)
		}

		note, err := notesRepo.Create(ctx, &entities.Note{
			Content: noteLine,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Note added:", note)
	},
}

func init() {
	rootCommand.AddCommand(addNoteCommand)

	addNoteCommand.PersistentFlags().String("n", "", "The note contents to add.")
}
