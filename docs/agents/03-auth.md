## Authentication & Security

- JWT access tokens (short-lived) and refresh-tokens
- Refresh tokens stored in Redis
- Passwords hashed using bcrypt
- Middleware must validate token type (access vs refresh)
