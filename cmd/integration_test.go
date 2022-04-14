package cmd

import (
	"bytes"
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"note-logger/internal/databases/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// for integration testing, let's start with a clean DB
	filename, err := sqlite.DBFilename()
	if err != nil {
		log.Fatalln(err)
	}

	os.Remove(filename)

	exitVal := m.Run()

	os.Exit(exitVal)
}

func runCommand(args []string) (string, error) {
	output := new(bytes.Buffer)

	rootCommand.SetOut(output)
	rootCommand.SetErr(output)

	rootCommand.SetArgs(args)
	err := rootCommand.Execute()

	return output.String(), err
}

var noteDetailsRegex = regexp.MustCompile(`^(\d*) - [\w ]*:\d*:\d*: (.*)`)
var noteDeletedRegex = regexp.MustCompile(`Note deleted.`)

func getNoteDetails(output string) ([]int, []string) {
	noteIDs := make([]int, 0)
	noteContents := make([]string, 0)

	actualLines := strings.Split(output, "\n")

	for _, line := range actualLines {
		matches := noteDetailsRegex.FindAllStringSubmatch(line, -1)

		if len(matches) > 0 {
			for i, match := range matches[0] {
				switch i {
				case 1:
					noteID, _ := strconv.Atoi(match)
					noteIDs = append(noteIDs, noteID)
				case 2:
					noteContents = append(noteContents, match)
				}
			}
		}
	}

	return noteIDs, noteContents
}

func TestIntegration(t *testing.T) {
	t.Run("error adding note without content", func(t *testing.T) {
		_, err := runCommand([]string{"add-note"})
		assert.Equal(t, errors.New("note content required"), err)
	})

	t.Run("add and then delete a note", func(t *testing.T) {
		actual, err := runCommand([]string{"add-note", "-c", "this is a new note"})
		assert.NoError(t, err)

		noteIDs, noteContents := getNoteDetails(actual)
		require.Equal(t, 1, len(noteIDs))
		require.Equal(t, 1, len(noteContents))

		assert.Equal(t, 1, noteIDs[0])
		assert.Equal(t, "this is a new note", noteContents[0])

		actual, err = runCommand([]string{"delete-note", "-i", strconv.Itoa(noteIDs[0])})
		assert.NoError(t, err)

		assert.Regexp(t, noteDeletedRegex, actual)
	})

	t.Run("adds a few notes and then lists them, then cleans up", func(t *testing.T) {
		_, err := runCommand([]string{"add-note", "-c", "note #1"})
		assert.NoError(t, err)

		_, err = runCommand([]string{"add-note", "-c", "note #2"})
		assert.NoError(t, err)

		_, err = runCommand([]string{"add-note", "-c", "note #3"})
		assert.NoError(t, err)

		actual, err := runCommand([]string{"list-notes", "-s", "10 minutes ago", "-e", "now"})
		assert.NoError(t, err)

		noteIDs, noteContents := getNoteDetails(actual)
		require.Equal(t, 3, len(noteIDs))
		require.Equal(t, 3, len(noteContents))

		assert.Equal(t, 1, noteIDs[0])
		assert.Equal(t, "note #1", noteContents[0])

		assert.Equal(t, 2, noteIDs[1])
		assert.Equal(t, "note #2", noteContents[1])

		assert.Equal(t, 3, noteIDs[2])
		assert.Equal(t, "note #3", noteContents[2])

		_, err = runCommand([]string{"delete-note", "-i", strconv.Itoa(noteIDs[0])})
		assert.NoError(t, err)

		_, err = runCommand([]string{"delete-note", "-i", strconv.Itoa(noteIDs[1])})
		assert.NoError(t, err)

		_, err = runCommand([]string{"delete-note", "-i", strconv.Itoa(noteIDs[2])})
		assert.NoError(t, err)
	})
}
