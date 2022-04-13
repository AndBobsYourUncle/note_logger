package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
	"github.com/tj/go-naturaldate"
)

var listNotesCommand = &cobra.Command{
	Use:   "list-notes",
	Short: "Lists the existing notes",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		beginningTimeString, err := cmd.Flags().GetString("start")
		if err != nil {
			log.Fatal(err)
		}

		if beginningTimeString == "" {
			log.Fatal(errors.New("beginning time required"))
		}

		endTimeString, err := cmd.Flags().GetString("end")
		if err != nil {
			log.Fatal(err)
		}

		if endTimeString == "" {
			log.Fatal(errors.New("end time required"))
		}

		beginningTime, err := naturaldate.Parse(beginningTimeString, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		endTime, err := naturaldate.Parse(endTimeString, time.Now())
		if err != nil {
			log.Fatal(err)
		}

		notesRepo, err := notes.NewRepository(&notes.Config{DB: sqliteDB})
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
	},
}

func init() {
	rootCommand.AddCommand(listNotesCommand)

	listNotesCommand.Flags().StringP("start", "s", "", "Start of the time window")
	listNotesCommand.Flags().StringP("end", "e", "", "End of the time window")
}
