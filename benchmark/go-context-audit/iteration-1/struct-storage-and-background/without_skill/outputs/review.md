# Go context review: `struct_storage.go`

**Task intent:** Review for request lifetime, context propagation, and hidden context misuse.

**Scope:** `auditfixture` package, file `struct_storage.go` only.

---

## Summary

The file combines two problems: **context is stored on a struct** (lifetime and API clarity issues) and **`LoadUser` runs the DB call under `context.Background()`**, so caller cancellation and deadlines never reach `QueryRowContext`. The constructor accepts a `context.Context` but that value is **never used** in the shown method, which can mislead readers into thinking the repository is request-scoped when the actual query is not.

---

## Findings

### 1. Context stored in `UserRepository` — **High**

| Item | Detail |
|------|--------|
| **Location** | Lines 8–10 (`UserRepository` fields), 13–17 (`NewUserRepository`) |
| **Issue** | `ctx` is held as a struct field. In Go, context is intended to be **per-operation** and passed into methods, not embedded in long-lived types. A repository often outlives a single HTTP/RPC request; a stored request context can become invalid, cancelled, or semantically wrong for later calls. |
| **Risk** | Wrong cancellation scope, accidental use of a stale context on a later call, or confusion about which operation the context applies to. |
| **Suggested fix** | Remove `ctx` from the struct. Pass `ctx context.Context` as the first parameter of each method that does I/O (e.g. `LoadUser(ctx context.Context, id string)`). Construct `UserRepository` with only dependencies like `db`. |

---

### 2. `LoadUser` uses `context.Background()` for the query — **High**

| Item | Detail |
|------|--------|
| **Location** | Lines 20–22 (`queryCtx := context.Background()`, `QueryRowContext`) |
| **Issue** | The database call explicitly opts out of the caller’s context. Timeouts and cancellation from the upstream request do not apply; work can continue after the client has gone away. |
| **Risk** | Resource leaks (connections, goroutines), slower failure under load, and violation of typical “request-bound work” expectations. |
| **Suggested fix** | Use the context passed into the method (after refactoring per finding 1), e.g. `return r.db.QueryRowContext(ctx, ...)`. If a subset of operations truly need a detached context, document why and consider a **bounded** timeout via `context.WithTimeout` from `context.Background()` only when justified—not as the default for request-path code. |

---

### 3. Unused / misleading stored context — **Medium**

| Item | Detail |
|------|--------|
| **Location** | Field `r.ctx` (lines 9, 15); contrast with `LoadUser` (lines 20–22) |
| **Issue** | `NewUserRepository` requires `ctx` and stores it, but `LoadUser` does not use `r.ctx`. This hides the fact that the repository does **not** propagate that context to I/O. |
| **Risk** | Maintainers may assume `r.ctx` is used for all operations and miss that queries are uncancellable. |
| **Suggested fix** | After removing struct-held context, the API matches behavior. If you must keep a “default” context for tests only, avoid production patterns that store request context on structs. |

---

## Risk legend

- **High:** Likely to cause incorrect cancellation/deadline behavior or encourage unsafe lifetime coupling.
- **Medium:** Misleading API or maintenance hazard that can lead to misuse.

---

## Positive notes

- `QueryRowContext` is used (the right API for cancellation-aware SQL), but it is wired to the wrong context in this file.

---

## Suggested shape (illustrative only)

```go
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) LoadUser(ctx context.Context, id string) error {
	return r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id).Err()
}
```

(Not applied to source; for review purposes only.)
