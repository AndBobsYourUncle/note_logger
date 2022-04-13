package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
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

		notesRepo, err := notes.NewRepository()
		if err != nil {
			log.Fatal(err)
		}

		err = notesRepo.Migrate(ctx)
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

		err = notesRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCommand.AddCommand(addNoteCommand)

	addNoteCommand.Flags().StringP("content", "c", "", "The note contents to add.")
}
