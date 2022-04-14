# Note Logger

## Summary
A very simple note-taking CLI you can use from the terminal that uses a SQLite DB to persist, and query, notes.

## Building/Installing
Ensure that you have Go 1.18 installed locally, and run:
```shell
go install
```

## Usage

On the first run, wherever the installed binary is located, a SQLite DB will be created called `notes.sqlite`.

This stores all the notes logged so far, and allows listing them back in the future.

If you have not compiled the app, but instead use `go run main.go note-logger`, then the DB will end up being located in a temp folder somewhere (differs by operating system, and Go installation).

### Add a Note

```shell
note-logger add-note -c "Some new note!"
```

This will result in some output:

```shell
Note added:
2 - Apr 12 16:32:31: Some new note!
```

You'll get the note's ID, the timestamp when it was created, and the contents.

### Delete a Note

```shell
note-logger delete-note -i 2
```

As an argument, you pass in the ID of the note to delete. This will be the output, if successful:

```shell
2022/04/12 16:33:46 Note deleted.
```

### List Notes

The listing of existing notes takes in two arguments, the beginning, and end for the time period. The values are interpreted using the available English-friendly values compatible with the [go-naturaldate](https://github.com/tj/go-naturaldate) package.

For example, here is listing all notes for today:

```shell
note-logger list-notes -s "beginning of today" -e "now"
```

Or even all of the week so far:

```shell
note-logger list-notes -s "beginning of week" -e "now"
```

You'll get output like this:

```shell
1 - Apr 12 16:26:19: First note with it all working!
2 - Apr 12 16:37:16: Another note for sample!
```

Similar to when you create a note, you'll get the note's ID, the timestamp, and the content. You can then retroactively delete notes this way using the `delete-note` command.

## Bash Functions

Executing the commands this way takes time, and perhaps it might be more convenient to type something simple into the terminal. Here are some sample Bash functions that you can add to your `.bashrc` file that make it easier to do common things:

```bash
function note() {
  quoted_note="$@"

  note-logger add-note -c $quoted_note
}

function delnote() {
  note-logger delete-note -i $@
}

function notes_today() {
  note-logger list-notes -s "beginning of today" -e "now"
}

function notes_week() {
  note-logger list-notes -s "beginning of week" -e "now"
}
```

With this, adding a note can be as simple as typing this in your terminal:
```bash
note Here is a new note!
```

## Contributing

Contributions are definitely welcome, so feel free to open a PR adding whatever new functionality you might like.

### DB Migrations

Database migrations are handled in the database migration wrapper here:
[sqlite.go](https://github.com/AndBobsYourUncle/note_logger/blob/master/internal/databases/sqlite/sqlite.go)

Adding a migration is as simple as adding an element to the migrations array:
```go
var migrations = []migration{
  {migrationName: "create notes table", migrationQuery: createTableIfNotExistsQuery},
  {migrationName: "add notes created_at index", migrationQuery: createIndexIfNotExistsQuery},
}
```

Migrations are run within a transaction, and any error results in a rollback. On any execution of the app, it checks the internal `pragma` setting, and if it is behind, it runs each migration in order to reach the required migration number.

Any migrations that might result in an error should get caught in the integration tests that run, as it starts with a fresh DB every time.

Here's a sample PR that adds a new index on the content column in the notes table:
https://github.com/AndBobsYourUncle/note_logger/pull/2
