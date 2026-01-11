# AGENTS.md

This file defines rules and expectations for AI agents working with this repository.

Отвечай всегда на русском языке.

---

## 1. Project Overview

Project is a backend service written in **Golang**.
Main stack:

- Go 1.25+
- PostgreSQL
- pgx / scany
- Fiber v3
- Redis (sessions, tokens)
- JWT-based authentication

Architecture:

- Clean / layered architecture
- Thin HTTP handlers
- Business logic in services
- Database access via repositories

---

## 2. Code Style & Conventions

### Go conventions

- Follow standard Go formatting (`gofmt`)
- Prefer explicit error handling
- Avoid panics in application code
- Do not ignore returned errors

### Naming

- Structs: `User`, `OrderService`
- Interfaces: `UserRepository`, `TokenService`
- Files: `snake_case.go`
- Packages: lowercase, short, meaningful

---

## 3. Project Structure - in file docs/agents/01-structure.md

## 4. Database rules - in file docs/agents/02-database.md

## 5. Authentifcation rules - in file docs/agents/03-auth.md

## 6. Others rules - in file docs/agents/04-others.md

## 6. What AI Agents SHOULD Do

- Propose minimal, incremental changes
- Follow existing architecture
- Ask before introducing new dependencies
- Provide runnable Go code
- Respect existing patterns and naming
- Answer in Russian if possible

## 7. What AI Agents MUST NOT Do

- Rewrite large parts of the codebase without request
- Introduce ORMs
- Change public APIs silently
- Add magic or overly abstract solutions
- Assume business logic without context

## 8. When in Doubt

If requirements are unclear:

- Ask clarifying questions
- Offer multiple options with trade-offs
- Prefer simplicity over cleverness
