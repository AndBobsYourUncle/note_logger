package cmd

import (
	"context"
	"fmt"
	"log"
	"note-logger/internal/repositories/notes"
	"time"

	"github.com/tj/go-naturaldate"

	"github.com/spf13/cobra"
)

var listNotesCommand = &cobra.Command{
	Use:   "list-notes",
	Short: "Lists the existing notes",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		beginningTimeString, err := cmd.Flags().GetString("b")
		if err != nil {
			log.Fatal(err)
		}

		endTimeString, err := cmd.Flags().GetString("e")
		if err != nil {
			log.Fatal(err)
		}

		beginningTime, err := naturaldate.Parse(beginningTimeString, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		endTime, err := naturaldate.Parse(endTimeString, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		notesRepo, err := notes.NewRepository()
		if err != nil {
			log.Fatal(err)
		}

		err = notesRepo.Migrate(ctx)
		if err != nil {
			log.Fatal(err)
		}

		notesRes, err := notesRepo.ListBetween(ctx, beginningTime, endTime)
		if err != nil {
			log.Fatal(err)
		}

		for _, note := range notesRes {
			fmt.Printf("%v - %v: %v\n", note.ID, note.CreatedAt.Format(time.Stamp), note.Content)
		}

		err = notesRepo.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCommand.AddCommand(listNotesCommand)

	listNotesCommand.PersistentFlags().String("b", "", "Beginning of the time window")
	listNotesCommand.PersistentFlags().String("e", "", "End of the time window")
}
