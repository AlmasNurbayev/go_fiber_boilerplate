### Configuration

Configuration via ENV variables
No hardcoded secrets
Use .env.example for documentation
Configuration loading must fail fast

### Error Handling

Use typed errors (var ErrNotFound = errors.New(...))
Map internal errors to HTTP errors in handlers only
Do not return raw DB errors to clients

### Testing Rules

Unit tests preferred over integration tests
Mock repositories in service tests
Use table-driven tests
Do not use t.Skip without explanation
