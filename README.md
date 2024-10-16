# Praja Organizations

PORGS—Praja Organizations is a web application for managing people, work, and discussions
within community and government organizations.

## Use

> 🚧🚧🚧 Under construction. Not ready for use. 🚧🚧🚧

## Contribute

> ⚠ Pull requests are not accepted at this time.

## Develop

[Setup](https://go.dev/doc/install) a Go development environment.

Clone the source code repository:
```shell
git clone https://github.com/praja-dev/porgs.git
```

Run the web app:
```shell
cd porgs
go run ./cmd/porgs
```

Open http://localhost:8642 on a web browser.

Use username `admin` and password `123` to sign-in. 

The home page (`/home`) lists the links to access functionality contributed by the active plugins.

Enter `Ctrl+C` to stop the porgs web app.

Run the web app again—this time loading example data from `./examples/lk/data` directory.
```shell
PORGS_LOAD_DIR=./examples/lk/data go run ./cmd/porgs
```

## Design

Overall design goals:
- **Simple to use**: The system should always present only what the user needs.
- **Simple to develop**: The entire system should be understandable by a single person.
- **Simple to maintain**: The system should deploy as a single binary with no external dependencies.
- **Simple to extend**: Features should be implemented as plugins atop a core system that handles essentials.
