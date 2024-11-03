# Praja Organizations

PORGS—Praja Organizations is a web application for managing people, work, and discussions
within community and government organizations.

## Use

> 🚧🚧🚧 Under construction. Not ready for use. 🚧🚧🚧

## Contribute

> ⚠ Pull requests are not accepted at this time.

## Develop

[Setup](https://go.dev/doc/install) a Go development environment.

On UNIX inspired systems:
```shell
mkdir ~/praja-dev
cd ~/praja-dev
git clone https://github.com/praja-dev/porgs.git
git clone https://github.com/praja-dev/lk-data.git
cd ~/praja-dev/porgs
go run ./cmd/porgs --load=~/praja-dev/lk-data/admin
```

On Windows:
```shell
mkdir %USERPROFILE%\praja-dev
cd %USERPROFILE%\praja-dev
git clone https://github.com/praja-dev/porgs.git
git clone https://github.com/praja-dev/lk-data.git
cd %USERPROFILE%\praja-dev\porgs
go run .\cmd\porgs --load=%USERPROFILE%\praja-dev\lk-data\admin
```

Open http://localhost:8642 on a web browser.

Use username `admin` and password `123` to sign-in. 

The home page (`/home`) lists the links to access functionality contributed by the active plugins.

Enter `Ctrl+C` to stop the porgs web app.

When starting again, don't use the `--load` argument.
```shell
go run ./cmd/porgs
```

The database is created in the same directory and has the name `porgs.db`.
This can be changed using the `PORGS_DSN` environment variable.

The database consists of three files:
- `porgs.db` - the main SQLite database file containing all tables, indexes etc.
- `porgs.db-wal` - write-ahead log file needed for WAL mode
- `porgs.db-shm` - shared memory file needed for WAL mode

You can safely delete these files to start fresh.

The [praja-dev/lk-data](https://github.com/praja-dev/lk-data.git) project is for creating a base dataset
for Sri Lanka starting from the country level (level 0 admin unit) down to the village level (level 4 admin unit).

For a truncated version of the lk-data dataset, refer to the [examples/lk/data directory](examples/lk/data) in the porgs project.  

## Design

Overall design goals:
- **Simple to use**: The system should always present only what the user needs.
- **Simple to develop**: The entire system should be understandable by a single person.
- **Simple to maintain**: The system should deploy as a single binary with no external dependencies.
- **Simple to extend**: Features should be implemented as plugins atop a core system that handles essentials.
  