## Database Rules

- PostgreSQL is the only DB
- Use `pgx` (no `database/sql`)
- SQL must be written explicitly (no ORM)
- Prefer `scany` for scanning rows
- returns only lib/dbErrors types
- Handle `NULL` values using `null.*` types

---
