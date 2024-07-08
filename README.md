# Praja Organizations

Praja Organizations (PORGS) is a software product for managing people, work,
and discussions within government, non-government, and community organizations.

## Use

> ⚠ PORGS is currently not ready for use as it is in early stage active development.

## Contribute

> ⚠ Pull requests are not accepted at this time.

## Develop

[Setup](https://go.dev/doc/install) a Go development environment.

Run the program
```shell
go run ./cmd/porgs
```

## Deploy

## Design

Overall design goals:
- **Simple to use**: The system should always present only what the user needs.
- **Simple to develop**: The entire system should be understandable by a single person.
- **Simple to maintain**: The system should deploy as a single binary with no external dependencies.
- **Simple to extend**: Features should be implemented as plugins atop a core system that handles essentials,
 integrating plugins for additional functionality.


## Product

**Background**: The initial requirements for the product are driven by the first client,
a large government organization in Sri Lanka. This organization features a significant
number of employees and a frequent rate of inter-branch transfers. However, PORGS is being
designed from the outset to be widely applicable to other government and community organizations.
