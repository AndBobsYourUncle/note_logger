package cmd

import (
	"context"
	"errors"
	"time"

	"note-logger/internal/databases/sqlite"
	"note-logger/internal/repositories/notes"

	"github.com/spf13/cobra"
	"github.com/tj/go-naturaldate"
)

var listNotesCommand = &cobra.Command{
	Use:   "list-notes",
	Short: "Lists the existing notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		beginningTimeString, err := cmd.Flags().GetString("start")
		if err != nil {
			return err
		}

		if beginningTimeString == "" {
			err := errors.New("beginning time required")
			return err
		}

		endTimeString, err := cmd.Flags().GetString("end")
		if err != nil {
			return err
		}

		if endTimeString == "" {
			err := errors.New("end time required")
			return err
		}

		beginningTime, err := naturaldate.Parse(beginningTimeString, time.Now())
		if err != nil {
			return err
		}

		endTime, err := naturaldate.Parse(endTimeString, time.Now())
		if err != nil {
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

		notesRes, err := notesRepo.ListBetween(ctx, beginningTime, endTime)
		if err != nil {
			return err
		}

		for _, note := range notesRes {
			cmd.Printf("%v - %v: %v\n", note.ID, note.CreatedAt.Format(time.Stamp), note.Content)
		}

		return nil
	},
}

func init() {
	rootCommand.AddCommand(listNotesCommand)

	listNotesCommand.Flags().StringP("start", "s", "", "Start of the time window")
	listNotesCommand.Flags().StringP("end", "e", "", "End of the time window")
}
