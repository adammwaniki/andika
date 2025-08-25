# Andika

==========================================================================================

Andika is an intuitive notes management service

## Tech Stack

This website is built using the GOTH stack:

- Go: an open-source, statically-typed and compiled language designed with systems programming in mind.
- Templ: an open-source HTML templating language for Go.
- HTMX: an open-source front-end JavaScript library that gives you access to AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML

## Andika API

Andika builds on a CLI version control system I had built [VCS CLI](https://github.com/adammwaniki/vcs-cli).
Below is a chart of the endpoints that exist in the current project and how they are mapped to the actions in the CLI.

### CLI <-> API Mapping

```json
| CLI Command                                   | HTTP Method     | Endpoint                      | Request Body                                   | Response                                                           |
| --------------------------------------------- | --------------- | ----------------------------- | ---------------------------------------------- | ------------------------------------------------------------------ |
| `create <file>`                               | `POST`          | `/notes/{noteName}`           | `{ "content": "string" }`                      | `{ "message": "note created", "hash": "..." }`                     |
| `view <file>`                                 | `GET`           | `/notes/{noteName}`           | –                                              | `{ "noteName": "...", "content": "..." }`                          |
| `edit <file>` (overwrite/append combined API) | `PUT`           | `/notes/{noteName}`           | `{ "mode": "append|overwrite", "content": "string" }` | `{ "message": "note updated", "hash": "..." }`             |
| `append <file>`                               | `PUT`           | `/notes/{noteName}`           | `{ "mode": "append", "content": "string" }`    | same as above                                                      |
| `overwrite <file>`                            | `PUT`           | `/notes/{noteName}`           | `{ "mode": "overwrite", "content": "string" }` | same as above                                                      |
| `list`                                        | `GET`           | `/notes`                      | –                                              | `[ "note1", "note2", ... ]`                                        |
| `list_snaps <file>`                           | `GET`           | `/notes/{noteName}/snapshots` | –                                              | `[ "snapshot1.gob", "snapshot2.gob", ... ]`                        |
| `snap view <snapshotHash>`                    | `GET`           | `/snapshots/{hash}`           | –                                              | `{ "noteName": "...", "content": "..." }`                          |
| `snap restore <snapshotHash>`                 | `POST`          | `/snapshots/{hash}/restore`   | –                                              | `{ "message": "restored note", "noteName": "...", "hash": "..." }` |
| `help`                                        | `GET`           | `/help`                       | –                                              | `{ "commands": [ ... ] }`                                          |

```

### Running Locally (Quickstart)

In your terminal, navigate to the vcs/ directory and run `air`

#### Create A New Note

Creates a new note with the provided content.

Method: `POST`
URL: `http://localhost:8160/notes/myNote`
Body:

```json
{
  "content": "Hello world, this is my first note"
}
```

Response (201 Created):

```json
{
  "message": "note created",
  "noteName": "myNote",
  "hash": "a7f3c9b7..." 
}
```

#### View A Note

User can view the contents of the given note

Method: `GET`
URL: `http://localhost:8160/notes/myNote`
Response (200 OK):

```json
{
  "noteName": "myNote",
  "content": "Hello world, this is my first note"
}
```

#### Edit A Note: Append mode

User can append to the given note

Method: `PUT`
URL: `http://localhost:8160/notes/myNote`
Body:

```json
{
  "mode": "append",
  "content": "\nAdding more content here"
}
```

Response:

```json
{
  "message": "note updated",
  "noteName": "myNote",
  "hash": "b9ac21f1..."
}
```

#### Edit A Note: Overwrite mode

User can overwrite the given note

Method: `PUT`
URL: `http://localhost:8160/notes/myNote`
Body:

```json
{
  "mode": "overwrite",
  "content": "This is new overwritten content"
}
```

Response:

```json
{
  "message": "note updated",
  "noteName": "myNote",
  "hash": "c2de34ab..."
}
```

#### List All Notes

User can list all notes

Method: `GET`
URL: `http://localhost:8160/notes`
Response (200 OK):

```json
[
  "myNote",
  "anotherNote"
]
```

#### List Snapshots of a Note

User can list all snapshots of a given note in chronologically descending order

Method: `GET`
URL: `http://localhost:8160/notes/myNote/snapshots`
Response (200 OK):

```json
[
  "a7f3c9b7a4.gob",
  "b9ac21f123.gob",
  "c2de34ab44.gob"
]
```

#### View Snapshot

User can view content of a given note snapshot

Method: `GET`
URL: `http://localhost:8160/snapshots/a7f3c9b7a4...`
Response (200 OK):

```json
{
  "noteName": "myNote",
  "content": "Hello world, this is my first note",
  "hash": "a7f3c9b7a4"
}
```

#### Restore Snapshot

User can restore a note to a given snapshot

Method: `POST`
URL: `http://localhost:8160/snapshots/a7f3c9b7a4/restore`
Response (200 OK):

```json
{
  "message": "note restored",
  "noteName": "myNote",
  "hash": "a7f3c9b7a4"
}
```

#### Help

Help menu for commands

Method: `GET`
URL: `http://localhost:8160/help`
Response (200 OK):

```json
{
  "commands": [
    "POST   /notes/{noteName}            -> create a note",
    "GET    /notes/{noteName}            -> view latest note content",
    "PUT    /notes/{noteName}            -> edit/append/overwrite note",
    "GET    /notes                       -> list all notes",
    "GET    /notes/{noteName}/snapshots  -> list all snapshots of a note",
    "GET    /snapshots/{hash}            -> view snapshot content",
    "POST   /snapshots/{hash}/restore    -> restore note to snapshot"
  ]
}
```

## Closing Remarks

Feel free to explore the site and see these principles in action!

Contributions and suggestions are welcome!
