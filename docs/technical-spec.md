## PORGS — Technical Specification

## Overview

**PORGS (Praja Organizations)** is a monolithic, single-repository web application designed for high portability and
decentralized integration. Built as a sovereign node, PORGS leverages a **Federated Block Architecture** to enable
independent instances to compose a resilient, nationwide information system.

- **Sovereign Deployment:** Each instance is a self-contained environment managing its own identities, hierarchies, and
  task lifecycles. It is designed for low-overhead deployment on local hardware or cloud infrastructure.
- **Federated Interoperability:** Instances establish secure, selective communication channels with other sovereign
  PORGS nodes. This allows for the organic emergence of hierarchical (Sub-Org) and horizontal (Peer-Org) relationships
  without a central point of failure.
- **Opinionated Data Primitives:** The system enforces a strict but flexible data model centered on **Identities**, *
  *Organizations**, and **Memberships**, ensuring data consistency and interoperability across the network as it scales.
- **Core Functionality & Extensibility:** PORGS provides an "essentials-only" subset of ERP, HR, and Project Management
  features, acting as a lightweight orchestration layer that can integrate with specialized external systems.

## Technical Choices

The technical choices are driven by these principles: **Simplicity**—ensuring a single engineer can comprehend the
entire system. **Self-containment**—prioritizing security and longevity by minimizing external dependencies. When a
dependency is unavoidable, we choose stable, mature tools that align with these values and leverage the full power of
the modern web platform.

- **[Kotlin Language](https://kotlinlang.org):** A modern, expressive JVM language. Fully Java-interoperable, granting
  access to the entire Java ecosystem without friction. Features like null safety, data classes, and coroutines
  eliminate boilerplate, allowing the entire codebase to remain "head-sized."
- **[Ktor Web Framework](https://ktor.io):** A lightweight, coroutine-based framework where behavior is explicit and
  composable. It avoids the "magic" of larger frameworks, ensuring there are no hidden runtime surprises.
- [SQLite](https://sqlite.org): Robust single-file in-process database. Used in WAL mode for one writer and multiple
  readers.
    - [xerial/sqlite-jdbc driver](https://github.com/xerial/sqlite-jdbc): Bundles native SQLite binaries into a
      single jar file. Adds about 10MB to the jar file (versus if only the binary for the target platform was bundled)
      but gives a seamless cross platform experience.
    - [Exposed](https://www.jetbrains.com/exposed/): Kotlin DSL for database access by JetBrains.
    - [Hikari Connection Pool](https://github.com/brettwooldridge/HikariCP): One pool for writes (maxConnections = 1)
      and one pool for
      reads used.
- **Netty Engine:** A battle-tested, non-blocking I/O engine. Embedded directly within the application, it eliminates
  the need for external server management.
- **Gradle with Kotlin DSL:** An industry-standard build tool using the same language as the application. Provides type
  safety and IDE autocompletion within build scripts, removing the context-switch to Groovy.
- **Eclipse Temurin (LTS):** A production-grade, openly licensed JDK distribution. Utilizing the same Long-Term Support
  release in both development and production eliminates environment-specific regressions.
- **YAML Configuration:** A human-readable, minimal syntax for application settings. Natively supported by Ktor, it
  allows for clear, version-controlled configuration without complex tooling.
- **[Thymeleaf Templates](https://www.thymeleaf.org):** A mature, server-side Java template engine that prioritizes
  the "Natural Templating" philosophy. While it sacrifices a fraction of raw execution speed compared to niche compiled
  engines, it provides a superior developer and designer experience. Templates remain valid HTML5 that can be opened and
  styled in any browser without a running server, ensuring the UI layer remains accessible to non-developers.
- **[Datastar](https://data-star.dev):** A modern, server-first reactivity library (~14kB) that replaces the HTMX +
  Alpine.js combination. It allows the server to drive state while the browser reacts, keeping application logic
  centralized and significantly reducing frontend surface area.
- **Plain CSS:** Leverages modern browser capabilities (Custom Properties, Grid, Flexbox, Cascade Layers) to build
  sophisticated UIs without a CSS framework. This avoids build steps and the inevitable "upgrade churn" of third-party
  stylesheets.
- **Plain JS:** Uses standard ES modules and modern APIs (`fetch`, `async/await`) for the minimal client-side logic not
  handled by Datastar. No bundlers or transpilers are required; the browser executes the code exactly as written.
- **Markdown Documentation:** Light-weight, version-controlled documentation that lives alongside the code. It is
  equally readable by humans, browsers, and LLM-assisted workflows.

## Configuration

Non-secret static application configuration is stored in [application.yaml](../src/main/resources/application.yaml) —
e.g. `app.name` and `app.version`. This is read once at startup into an `AppConfig` data class and stored as an
application-scoped `AttributeKey`. Sensitive application onfiguration is stored in environment variables — e.g. security
keys.

## Build

Gradle with the Kotlin DSL is used as the build tool.

### Dependencies

A Version Catalog — `gradle/libs.versions.toml` — is used to manage all Gradle plugin and library versions in one place.
Dependencies are declared in `build.gradle.kts` referencing the catalog via `libs.*` accessors.

Constraints are used in `build.gradle.kts` to force specific transitive dependency versions (e.g. the latest Netty) when
upstream libraries have not yet updated.

Dependency lines in Gradle scripts and in source code are always alphabetically sorted for consistency.

### Build Configuration

Gradle is configured across these files:

- `settings.gradle.kts` — the entry point for any Gradle build. Gradle reads this file first, before `build.gradle.kts`.
  It establishes the project name (which becomes the artifact name), and in our case also configures
  `dependencyResolutionManagement` to declare `mavenCentral()` as the sole repository for resolving dependencies.
  Centralising repository configuration here rather than in `build.gradle.kts` is the modern Gradle convention for
  single-project builds.
- `gradle/libs.versions.toml` — version catalog; single source of truth for all dependency and plugin versions
- `build.gradle.kts` — applies plugins, sets the main class, configures the JVM toolchain, and declares dependencies
- `gradle/wrapper/gradle-wrapper.properties` — pins the Gradle wrapper version to ensure consistent builds across
  environments

### Tasks

**`./gradlew run`** starts the embedded Netty server and begins serving requests on the configured port. This is the
primary command during development — it compiles the application and runs it in the same JVM process.

**`./gradlew test`** compiles and runs the test suite using the Ktor test engine, which spins up the application
in-process without binding to a real port. Run this before committing any change.

**`./gradlew build`** is the full build pipeline: compiles sources, runs tests, and assembles the output artifacts. It
is what CI runs. It does not produce a deployment-ready artifact on its own — use one of the tasks below for that.

**`./gradlew buildFatJar`** produces a single self-contained JAR under `build/libs/` that includes the application code
and all its dependencies. This is the simplest deployment artifact — copy it to a server with a JRE installed and run
`java -jar porgs-all.jar`. No external application server or container runtime required.

**`./gradlew buildImage`** uses the Ktor Gradle plugin to build an OCI-compliant container image from the fat JAR. Use
this when deploying to a container runtime (Docker, Podman) or an orchestrator (Kubernetes). The image bundles the JRE,
making it fully self-contained and portable across any OCI-compatible host.

## Ktor

Understanding how a Ktor application boots requires understanding three distinct layers: the engine, the application,
and modules. Here is how they fit together.

**The Engine:** Ktor is not tied to a specific web server — it supports Netty, Jetty, and others as pluggable engines.
We use Netty. Rather than writing a `main` function that manually constructs and starts a server, we delegate to
`io.ktor.server.netty.EngineMain`:

```kotlin
fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}
```

`EngineMain` handles the complete server lifecycle: reading configuration, starting Netty, managing the coroutine
dispatcher, handling graceful shutdown on SIGTERM, and loading application modules. We get all of this for free.

`EngineMain` reads `src/main/resources/application.yaml` on startup. The key property is:

```yaml
ktor:
  application:
    modules:
      - dev.praja.porgs.ApplicationKt.bootstrap
```

This tells the engine to reflectively invoke a function named `bootstrap` defined in the `ApplicationKt` class (Kotlin
compiles top-level functions in `Application.kt` into a class named `ApplicationKt`). This decoupling means the engine
has no compile-time dependency on the application logic — you could swap engines without changing a line of application
code.

**Kotlin Extension Functions: ** A Kotlin extension function lets you add a function to an existing class without
modifying or subclassing it. The syntax `fun Application.bootstrap()` defines a function that appears to be a method on
the `Application` class (which comes from the Ktor library), but is actually declared in our own code. Inside that
function, `this` refers to the `Application` instance, giving access to everything on it — configuration, attributes,
the logger, the coroutine scope.

This is the central pattern Ktor builds on. Every piece of application setup — installing plugins, defining routes,
connecting to a database — is written as an extension function on `Application`. This prevents a single "God Object"
configuration file from growing unboundedly. Each concern lives in its own file, each expressed as a focused extension
function.

**Modules:** In Ktor terminology, a *module* is any `fun Application.xxx()` extension function registered in
`application.yaml` - in this case `bootstrap` — which calls further extension functions:

```kotlin
fun Application.bootstrap() {
    // Read config, store in attributes, install plugins, ...
    configureRouting()   // extension function defined in Routing.kt
}
```

`configureRouting()` is itself an extension function on `Application`, defined in `Routing.kt`. Calling it from
`bootstrap()` works because both share the same `Application` receiver — Kotlin resolves the call on `this`.

**Plugins:** Ktor features (middleware, in Ktor terms) are provided as *plugins*. You install a plugin with
`install(PluginName) { /* config block */ }`inside a module. Plugins intercept the request/response pipeline. Plugins
are composable and order-dependent: i.e., they are applied in the order they are installed.

**Application Attributes:** `Application` carries an `attributes` map (`AttributeKey<T>` typed) for storing
application-scoped singletons — things that are created once at startup and shared across the lifetime of the server. We
use this to store `AppConfig` so that any extension function with access to the `Application` receiver (which is all of
them) can retrieve it without passing it around explicitly:

```kotlin
attributes.put(AppConfigKey, appConfig)   // Set in bootstrap()
val appConfig = attributes[AppConfigKey]  // Get in configureRouting()
```

**Summary of the call chain at startup**

```
JVM starts
  → main() in Application.kt
    → EngineMain reads application.yaml
      → Netty starts, binds to port 8080
        → EngineMain reflectively calls Application.bootstrap()
          → AppConfig read from config, stored in attributes
          → configureRouting() called → routes registered
            → server ready to accept requests
```

## Logging

**The two-layer model**

Java logging has historically used a facade/implementation split. The facade — **SLF4J** (Simple Logging Facade for
Java) — is what application code calls. It defines a standard `Logger` interface and is a compile-time dependency only.
The implementation — **Logback** — is what actually writes log output. It is a runtime dependency.

This separation means you can swap the logging implementation without changing any application code. Ktor itself logs
via SLF4J for the same reason.

**Configuration**

Logback is configured in `src/main/resources/logback.xml`. The current configuration writes all `INFO`-and-above log
output to stdout in a timestamped format:

```xml

<pattern>%d{YYYY-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n</pattern>
```

The `%logger{36}` token prints the logger name (typically the fully-qualified class name) truncated to 36 characters,
which is enough to identify the source.

**Log levels**

SLF4J defines five levels, in ascending severity:

- `TRACE` — fine-grained diagnostic detail; almost never needed in normal operation
- `DEBUG` — diagnostic information useful during development
- `INFO` — normal operational events (server started, request handled)
- `WARN` — unexpected but recoverable situations
- `ERROR` — failures that require attention

Setting a level on a logger means that level *and above* are emitted. Setting `INFO` suppresses `DEBUG` and `TRACE`.

**Changing the root log level**

The root logger in `logback.xml` controls the default level for everything:

```xml

<root level="INFO">
    <appender-ref ref="STDOUT"/>
</root>
```

Change `INFO` to `DEBUG` to see debug output from Ktor internals and application code. Change to `WARN` to suppress
informational noise in production.

**Changing the level for a specific package or class**

Add a `<logger>` element targeting the package or fully qualified class name:

```xml
<!-- Suppress verbose Netty internal logging -->
<logger name="io.netty" level="WARN"/>

        <!-- Enable DEBUG only for our code -->
<logger name="dev.praja.porgs" level="DEBUG"/>
```

The most specific matching logger wins. This lets you dial up detail for one package while keeping the rest quiet.

**Setting the level at runtime via system property**

When running with `./gradlew run`, you can override the root level without editing `logback.xml` by passing a system
property:

```shell
./gradlew run -Dlogback.configurationFile=logback-debug.xml
```

Or maintain a separate `src/main/resources/logback-debug.xml` with `<root level="DEBUG">` and point to it when needed.
This keeps the debug configuration out of the production config file.

**Writing log output in application code**

Ktor injects a logger into the `Application` scope accessible as `log`:

```kotlin
fun Application.bootstrap() {
    log.info("Starting PORGS application")
    log.debug("AppConfig loaded: {}", appConfig)
}
```

The `{}` placeholder syntax is an SLF4J convention — it defers `toString()` evaluation until the message is actually
going to be emitted, avoiding string construction overhead at suppressed log levels. Prefer it over string interpolation
in log statements.

## Routing

Routing is covered at the architectural level in the *Entry Point and Modules* section above — specifically how
`configureRouting()` is an extension function on `Application`. This section covers the routing DSL itself.

**The routing block**

`routing { }` is a Ktor plugin that registers a route-matching tree on the application. Everything inside the block is
declarative — you describe the shape of URLs and what to do when they match:

```kotlin
fun Application.configureRouting() {
    routing {
        get("/") {
            call.respondTemplate("index", mapOf("appConfig" to appConfig))
        }
    }
}
```

**Route handlers**

Each route handler is a suspending lambda with a `PipelineContext<Unit, ApplicationCall>` receiver, but in practice you
only interact with `call` — the `ApplicationCall` representing the current request/response pair.
**HTTP method functions**

Ktor provides a function for each HTTP method:

```kotlin
get("/tasks") { }
post("/tasks") { }
put("/tasks/{id}") { }
delete("/tasks/{id}") { }
patch("/tasks/{id}") { }
```

**Path parameters**

Curly brace segments are captured as named parameters:

```kotlin
get("/tasks/{id}") {
    val id = call.parameters["id"]
        ?: return@get call.respond(HttpStatusCode.BadRequest)
    // use id
}
```

The `?:` (Elvis operator) handles the case where the parameter is missing. `return@get` is a labelled return — it exits
the `get`lambda early, a common Kotlin pattern for early returns from lambdas.

**Route grouping**

Routes can be nested under a common path prefix using `route()`:

```kotlin
routing {
    route("/tasks") {
        get { /* GET /tasks */ }
        post { /* POST /tasks */ }
        route("/{id}") {
            get { /* GET /tasks/{id} */ }
            put { /* PUT /tasks/{id} */ }
            delete { /* DELETE /tasks/{id} */ }
        }
    }
}
```

This mirrors the REST resource structure and avoids repeating path segments.

**Organising routes across files**

As the application grows, all routes should not live in a single `routing { }` block. The idiomatic approach is to
define further extension functions for each resource area and call them from `configureRouting()`:

```kotlin
// Routing.kt
fun Application.configureRouting() {
    routing {
        orgRoutes()
        taskRoutes()
    }
}

// TaskRoutes.kt
fun Route.taskRoutes() {
    route("/tasks") {
        get { }
        post { }
    }
}
```

Note the receiver type changes to `Route` (not `Application`) for these nested functions — `Route` is what you have
inside a `routing { }` block. `routing { }` itself returns a `Route`, so any extension function on `Route` can be called
inside it.

## Auth

TODO LATER.

## Localization

TODO LATER.

## User Interface

### Assets

Static assets (CSS, JavaScript, images, etc.) are served directly by Ktor using the `staticResources()` DSL, which
resolves files from the application's classpath. This ensures assets are bundled into the fat JAR and available in
production without any filesystem path configuration.

```shell
staticResources("/assets", "assets")
```

With this, the content of `src/main/resources/assets` is served at `/assets`.

Assets are organized into the following directories:

```
src/main/resources/
    assets/       ← Assets other than CSS, JS, and images
    css/          ← Plain CSS files
    images/       ← Image files
    javascript/   ← Plain JS files
```

Browser **cache busting** is handled via a version query parameter appended to asset URLs.

```html

<link rel="stylesheet" th:href="|/css/index.css?v=${appConfig.version}|">
```

This was chosen over asset fingerprinting due to simplicity.

### Templates

Thymeleaf is configured in `bootstrap()` using a `ClassLoaderTemplateResolver`. Templates live under
`src/main/resources/templates/` and are addressed by their path relative to that directory, without the `.html` suffix:

A controller renders a template by calling `call.respondTemplate()` with the template name and a model map:

```kotlin
call.respondTemplate(
    "home/index",
    mapOf("appConfig" to appConfig)
)
```

The model map's keys become available as Thymeleaf expression variables in the template (e.g. `${appConfig.version}`).

An ordinary template named `layout.html` defines fragments named `c_head` and `c_body` which holds the conent for the
`<head>` and `<body>` elements of a common web page.

Templates are organized by controller/feature area:

```shell
src/main/resources/templates/
  layout.html          ← The common layout
  fragments.html       ← Fragments for use in other templates — Datastar patch targets, reusable snippets
  home/
      index.html       ← GET /
  sandbox/
      index.html       ← GET /sandbox
```

### Reactivity

[Datastar](https://data-star.dev/) unifies backend-driven DOM updates (like htmx) and frontend signal-based reactivity (
like Alpine.js) into a single ~14 kB module. It is served from the application's own static resource path (
`/javascript/lib/datastar.js`) rather than a CDN, keeping the application self-contained and offline-capable.

The library is manually downloaded when a new version is released.

```shell
curl -L "https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.0-RC.8/bundles/datastar.js" \
  -o src/main/resources/javascript/lib/datastar.js
```

To get rid of the "Namespace errors for data-*:* attributes" in IntelliJ IDEA:

- Disable: Settings > Editor > Inspections > XML > Unbound namespace prefix

#### How it works

Datastar operates on two primitives:

- **Signals** — reactive client-side state declared in HTML via `data-signals`. All signals are automatically sent to
  the server as a JSON body on every backend request.
- **Attributes** — `data-*` directives that bind signals to the DOM (`data-text`, `data-show`, `data-class`, etc.) and
  wire up actions (`data-on:click`, `data-on:load`, etc.).

The server responds with one of two content types:

- `text/html` — the response is morphed directly into the DOM (simple cases).
- `text/event-stream` — zero or more SSE events are streamed, each patching signals or elements independently.

The backend is the single source of truth. The browser reacts.

#### SSE events

Ktor endpoints that serve Datastar responses set `Content-Type: text/event-stream` and write raw SSE using
`respondTextWriter`:

```kotlin
call.response.header(HttpHeaders.CacheControl, "no-cache")
call.respondTextWriter(contentType = ContentType.Text.EventStream) {
    write("event: datastar-patch-signals\n")
    write("data: signals {count: $next}\n")
    write("\n")
    flush()
}
```

Two event types are used:

- `datastar-patch-signals` — merges a partial signal object into the client signal store; any `data-text` or other
  bindings update automatically.
- `datastar-patch-elements` — morphs one or more HTML fragments into the DOM by matching element IDs.

#### Reading signals from the request

On `@post()` actions, Datastar sends all current signals as a flat JSON body. Parse them directly — there is no wrapper
key:

```kotlin
val body = call.receiveText()
val json = Json.parseToJsonElement(body).jsonObject
val count = json["count"]?.jsonPrimitive?.int ?: 0
```

On `@get()` actions, signals are sent as a `datastar` query parameter instead.

#### Template wiring

Signals are initialized on a container element and consumed by child elements:

```html

<div data-signals="{count: 0}">
    <div id="ds-count" data-text="'Count: ' + $count"></div>
    <button data-on:click="@post('/counter/increment')">Increment</button>
</div>
```

The `id` on patched elements is required — `datastar-patch-elements` morphs by ID.

## Database

The databas is a single file — `storage/porgs.db`, configured via the `app.db` key in `application.yaml`. The `storage/`
directory is created automatically on first boot if absent.

### Connection Pools

Two distinct HikariCP pools target the same `.db` file:

- **Writer pool** — `maximumPoolSize = 1`, plain JDBC URL (`jdbc:sqlite:<path>`)
- **Reader pool** — `maximumPoolSize = 10`, URI-mode URL with `mode=ro` (`jdbc:sqlite:file:<path>?mode=ro&uri=true`)

Read-only enforcement is applied at the SQLite driver level via the URI flag, not via HikariCP's `isReadOnly`, because
SQLite requires the read-only flag to be set before the connection is opened.

Both pools apply the following pragmas via `connectionInitSql` on every new connection:

- `PRAGMA journal_mode` — `WAL` mode enables concurrent reads during a write operation.
- `PRAGMA foreign_keys` — `ON` enforces referential integrity.
- `PRAGMA busy_timeout` — prevents immediate failure under write-contention.

### Concurrency & Dispatchers

All database work is dispatched to `Dispatchers.IO` to prevent blocking Ktor's event loop. Parallelism is capped to
mirror the pool sizes:

```kotlin
val readDispatcher = Dispatchers.IO.limitedParallelism(10)
val writeDispatcher = Dispatchers.IO.limitedParallelism(1)
```

Two top-level suspend functions provide the public interface for all database access:

```kotlin
suspend fun <T> dbRead(block: Transaction.() -> T): T
suspend fun <T> dbWrite(block: Transaction.() -> T): T
```

### Schema Definition

Table objects are defined in `src/main/kotlin/tables/` using Exposed's `LongIdTable`. This maps Kotlin `Long` to
SQLite's `INTEGER PRIMARY KEY`, which is an alias for the internal 64-bit `rowid`.

```kotlin
object People : LongIdTable("people") {
    val name = varchar("name", 255)
}
```

### Query Layer (Active Record Style)

Rather than a repository layer, queries are defined as `suspend` functions directly on each model's `companion object`.
This mirrors the ActiveRecord calling convention (`Person.find(1)`, `Person.all()`) while remaining idiomatic Kotlin.
The companion object functions call `dbRead` or `dbWrite` and map `ResultRow` values to the data class within the
transaction, ensuring connections are returned to the pool promptly.

```kotlin
data class Person(val id: Long, val name: String) {
    companion object {
        suspend fun find(id: Long): Person? = dbRead { }
        suspend fun all(): List<Person> = dbRead { }
    }
}
```

Model files are located in `src/main/kotlin/models/`. Each model owns its own query functions. Table objects in
`tables/` are an implementation detail of the model file and are not referenced outside it.

### Migration Workflow

Migrations are plain SQL files in `db/migrate/`. During initial development (pre-0.1.0), a single `_schema.sql` file
holds the full baseline schema and is updated in place.

From version 0.1.0 onward, migrations will follow a timestamp-prefixed naming convention and applied migration versions
will be tracked in a `schema_migrations` table in the database itself. Then, the application will refuse to boot if
unapplied migrations are detected.

### Performance Guidelines

- **Keep transactions small.** Do not perform network calls or heavy computation inside a `dbRead` or `dbWrite` block.
- **Map inside the transaction.** Convert `ResultRow` to model instances within the transaction block so the connection
  is released before the data is used.
- **Datastar integration.** Fetch data in a read transaction and render the Thymeleaf template outside of it.

## Data Model

TODO LATER.

## Features

TODO LATER.
