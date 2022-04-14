package cmd

import (
	"bytes"
	"errors"
	"log"
	"note-logger/internal/databases/sqlite"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

var noteAddedRegex = regexp.MustCompile(`^(\d*) - [\w ]*:\d*:\d*: (.*)`)
var noteDeletedRegex = regexp.MustCompile(`Note deleted.`)

func getNoteCreatedDetails(t *testing.T, output string) (int, string) {
	actualLines := strings.Split(output, "\n")

	assert.Equal(t, 3, len(actualLines))
	assert.Equal(t, "", actualLines[len(actualLines)-1])

	matches := noteAddedRegex.FindAllStringSubmatch(actualLines[1], -1)

	assert.Equal(t, 1, len(matches))
	assert.Equal(t, 3, len(matches[0]))

	noteID, err := strconv.Atoi(matches[0][1])
	assert.NoError(t, err)

	return noteID, matches[0][2]
}

func TestIntegration(t *testing.T) {
	t.Run("error adding note without content", func(t *testing.T) {
		_, err := runCommand([]string{"add-note"})
		assert.Equal(t, errors.New("note content required"), err)
	})

	t.Run("add and then delete a note", func(t *testing.T) {
		actual, err := runCommand([]string{"add-note", "-c", "this is a new note"})
		assert.NoError(t, err)

		noteID, contents := getNoteCreatedDetails(t, actual)

		assert.Equal(t, 1, noteID)
		assert.Equal(t, "this is a new note", contents)

		actual, err = runCommand([]string{"delete-note", "-i", strconv.Itoa(noteID)})
		assert.NoError(t, err)

		assert.Regexp(t, noteDeletedRegex, actual)
	})
}
