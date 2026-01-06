## Structure

-cmd

- server (entrypoint http-server)
- migrator (app for apply migrations to DB)
- seeder (app for create minimal data in DB)

- internal

  - app (temporary empty)
    -- config (code for load config/env)
  - db
    -- storage (init, close, methods for postgres)
  - httpApp (http-server, routes)
    -- dto (requests, responses struct types)
    -- handlers (http handlers)
    -- services (main business logic)
    -- middleware (prometheus, auth and others middlewares)
  - /service
  - lib (utilities - jwt, errors handling, validate etc.)
  - logger (initial slog)
  - models (DB tables struct types)

- migrtations (sql-files for apply changes up/down)
- seeds (sql-files for create minimal data in DB)

Rules:

- `cmd` contains only bootstrap code
- Cross-layer imports are forbidden (e.g. handler → repository)
- imports between different layers (storage, handlers, services) must be done through interfaces
- the storage layer uses its own types (models), service/handlers layers should use Dto
- To convert between models/dto services should use the github.com/jinzhu/copier package

---

6. Configuration

Configuration via ENV variables

No hardcoded secrets

Use .env.example for documentation

Configuration loading must fail fast

7. Error Handling

Use typed errors (var ErrNotFound = errors.New(...))

Map internal errors to HTTP errors in handlers only

Do not return raw DB errors to clients

8. Testing Rules

Unit tests preferred over integration tests

Mock repositories in service tests

Use table-driven tests

Do not use t.Skip without explanation
