# PORGS

PORGS—Praja Organizations is a web application for managing people, work, and discussions
within community and government organizations.

## Contribute

> ⚠ Pull requests are not accepted at this time.

## Use

> 🚧🚧🚧 Under construction. Not ready for use. 🚧🚧🚧

## Setup

Prerequisites:

- JDK — the latest LTS version
- Kotlin — the latest stable version

First, install [mise](https://mise.jdx.dev/installing-mise.html) — a development environment setup tool.

Then, install Java

```shell
mise use --global java@temurin-25
```

Install Kotlin

```shell
mise use --global kotlin@2.3.20
```

## Run

Run the server

```shell
./gradlew run
```

Then, open the default web browser at the server's root URL

```shell
open http://localhost:8080 
```

Clean and do a full build:

```shell
./gradlew clean
./gradlew build
```

## Deploy


Build a fat JAR that includes all that is needed:

```shell
./gradlew buildfatJar
```

Run fat JAR:

```shell
java -jar build/libs/porgs-all.jar
```

Build a Docker image

```shell
./gradlew buildImage
```

