# Andika

==========================================================================================

Andika is an intuitive notes management service

## Tech Stack

This website is built using the GOTH stack:

- Go: an open-source, statically-typed and compiled language designed with systems programming in mind.
- Templ: an open-source HTML templating language for Go.
- HTMX: an open-source front-end JavaScript library that gives you access to AJAX, CSS Transitions, WebSockets and Server Sent Events directly in HTML

## Andika API

Andika builds on a CLI version control system I had built accessible at [VCS CLI](https://github.com/adammwaniki/vcs-cli).
Below is a chart of the endpoints that exist in the current project and how they are mapped to the actions in the CLI.

### CLI <-> API Mapping

| CLI Command                                   | HTTP Method | Endpoint                                    | Request Body                                        | Response                                                                   |
| --------------------------------------------- | ----------- | ------------------------------------------- | --------------------------------------------------- | -------------------------------------------------------------------------- |
| `create <file>`                               | `POST`      | `/api/v1/notes`                             | `{ "name": "string", "content": "string" }`        | `{ "message": "note created", "id": "uuid", "name": "...", "hash": "..." }` |
| `view <file>`                                 | `GET`       | `/api/v1/notes/{id}`                        | –                                                   | `{ "id": "...", "name": "...", "content": "..." }`                        |
| `edit <file>` (overwrite/append combined API) | `PUT`      | `/api/v1/notes/{id}`                        | `{ "mode": "append|overwrite|edit", "content": "string" }` | `{ "message": "note updated", "id": "...", "name": "...", "hash": "..." }` |
| `append <file>`                               | `PUT`       | `/api/v1/notes/{id}`                        | `{ "mode": "append", "content": "string" }`        | same as above                                                              |
| `overwrite <file>`                            | `PUT`       | `/api/v1/notes/{id}`                        | `{ "mode": "overwrite", "content": "string" }`     | same as above                                                              |
| `list`                                        | `GET`       | `/api/v1/notes`                             | –                                                   | `[{ "id": "uuid", "name": "..." }, ...]`                                  |
| `list_snaps <file>`                           | `GET`       | `/api/v1/notes/{id}/snapshots`              | –                                                   | `{ "id": "...", "name": "...", "hashes": ["hash1", "hash2", ...] }`       |
| `snap view <snapshotHash>`                    | `GET`       | `/api/v1/notes/{id}/snapshots/{hash}`       | –                                                   | `{ "id": "...", "name": "...", "content": "...", "hash": "..." }`         |
| `snap restore <snapshotHash>`                 | `POST`      | `/api/v1/notes/{id}/snapshots/{hash}/restore` | –                                                 | `{ "message": "note restored", "id": "...", "name": "...", "hash": "..." }` |
| `delete <file>`                               | `DELETE`    | `/api/v1/notes/{id}`                        | –                                                   | `{ "message": "note deleted successfully", "id": "...", "name": "..." }`  |
| `help`                                        | `GET`       | `/api/v1/help`                              | –                                                   | `{ "message": "...", "commands": [...] }`                                 |

### Running Locally (Quickstart)

In your terminal, navigate to the backend/ directory and run `go run main.go` or use `air` for hot reloading.

The API will be available at `http://localhost:8160`

#### Create A New Note

Creates a new note with the provided name and content. The system generates a **unique UUID** for each note.

Method: `POST`

URL: `http://localhost:8160/api/v1/notes`

Body:

```json
{
  "name": "My First Note",
  "content": "Hello world, this is my first note"
}
```

Response (201 Created):

```json
{
  "message": "note created",
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "hash": "a7f3c9b7a4..."
}
```

#### View A Note

User can view the contents of the given note by its ID

Method: `GET`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496`

Response (200 OK):

```json
{
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "content": "Hello world, this is my first note"
}
```

#### Edit A Note: Append mode

User can append to the given note

Method: `PUT`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496`

Body:

```json
{
  "mode": "append",
  "content": "\nAdding more content here"
}
```

Response (200 OK):

```json
{
  "message": "note updated",
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "hash": "b9ac21f1..."
}
```

#### Edit A Note: Overwrite mode

User can overwrite the given note

Method: `PUT`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496`

Body:

```json
{
  "mode": "overwrite",
  "content": "This is new overwritten content"
}
```

Response (200 OK):

```json
{
  "message": "note updated",
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "hash": "c2de34ab..."
}
```

#### List All Notes

User can list all notes with their IDs and names

Method: `GET`

URL: `http://localhost:8160/api/v1/notes`

Response (200 OK):

```json
[
  {
    "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
    "name": "My First Note"
  },
  {
    "id": "c9425d0c-7325-4a26-b292-81f0bda8cb82",
    "name": "Another Note"
  }
]
```

#### List Snapshots of a Note

User can list all snapshots of a given note in chronological order

Method: `GET`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496/snapshots`

Response (200 OK):

```json
{
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "hashes": [
    "a7f3c9b7a4...",
    "b9ac21f1...",
    "c2de34ab..."
  ]
}
```

#### View Snapshot

User can view content of a specific note snapshot

Method: `GET`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496/snapshots/a7f3c9b7a4...`

Response (200 OK):

```json
{
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "content": "Hello world, this is my first note",
  "hash": "a7f3c9b7a4..."
}
```

#### Restore Snapshot

User can restore a note to a given snapshot

Method: `POST`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496/snapshots/a7f3c9b7a4.../restore`

Response (200 OK):

```json
{
  "message": "note restored",
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note",
  "hash": "a7f3c9b7a4..."
}
```

#### Delete A Note

User can delete a note and all its snapshots

Method: `DELETE`

URL: `http://localhost:8160/api/v1/notes/915dd725-21c1-4ab2-beff-e49e51d38496`

Response (200 OK):

```json
{
  "message": "note deleted successfully",
  "id": "915dd725-21c1-4ab2-beff-e49e51d38496",
  "name": "My First Note"
}
```

#### Help

Help menu for commands

Method: `GET`

URL: `http://localhost:8160/api/v1/help`

Response (200 OK):

```json
{
  "message": "Andika API v1 Help",
  "commands": [
    "POST   /notes                         -> create a note (body: {name, content})",
    "GET    /notes                         -> list all notes",
    "GET    /notes/{id}                     -> view latest content of a note",
    "PUT    /notes/{id}                     -> edit/append/overwrite note (body: {mode, content})",
    "DELETE /notes/{id}                     -> delete a note",
    "GET    /notes/{id}/snapshots           -> list all snapshots of a note",
    "GET    /notes/{id}/snapshots/{sid}     -> view a specific snapshot",
    "POST   /notes/{id}/snapshots/{sid}/restore -> restore note to a snapshot",
    "GET    /help                           -> display this help message"
  ]
}
```

## Key Differences from Original Design

- **UUID-based**: Notes are identified by UUIDs instead of names for better uniqueness
- **Versioned API**: All endpoints are prefixed with `/api/v1/` for API versioning
- **Structured Responses**: More consistent JSON response format with metadata
- **Delete Functionality**: Added ability to delete notes completely
- **Better Error Handling**: Comprehensive error responses with appropriate HTTP status codes

## Closing Remarks

Feel free to explore the API and see these principles in action!

Contributions and suggestions are welcome!
