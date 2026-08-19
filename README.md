# go-reco

Periodically fetches users and projects from the Asana API and saves each one
as its own JSON file under `data/`.

## Prerequisites

- Go 1.26+
- An Asana personal access token
- The Asana workspace GID you want to pull data from

## Setup

Clone the repo, then install dependencies:

```sh
go mod download
```

## Configuration

Copy the example env file and fill in your own values:

```sh
cp .env.example .env
```

`.env` variables:

| Variable         | Description                                                              |
|------------------|---------------------------------------------------------------------------|
| `ACCESS_TOKEN`   | Your Asana personal access token. Sent as `Authorization: Bearer <token>`. |
| `WORKSPACE_GID`  | The GID of the Asana workspace to fetch users/projects from.             |
| `PAGE_SIZE`      | Number of records requested per API page (Asana's `limit` query param).  |
| `FETCHING_MODE`  | How often to refetch. See table below.                                   |

`FETCHING_MODE` values:

| Value | Refetch interval | Notes                          |
|-------|-------------------|--------------------------------|
| `1`   | every 30 seconds  | default / fallback for any other value |
| `2`   | every 5 minutes   |                                 |

## Running

```sh
go run .
```

This loads `.env`, then fetches users and projects on a loop at the
configured interval. Output is written to:

- `data/users/<gid>.json`
- `data/projects/<gid>.json`

Stop the process with `Ctrl+C`.

## Running the tests

```sh
go test ./...
```
