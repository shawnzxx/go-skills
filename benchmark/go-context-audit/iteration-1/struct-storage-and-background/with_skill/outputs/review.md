# Go Context Audit: `struct_storage.go`

**Task intent:** Review Go code for context issues. Focus on request lifetime, context propagation, and whether this file hides any context misuse.

**Scope:** Single file `evals/files/struct_storage.go` (package `auditfixture`). Transitive callers/callees were not inspected.

---

## Summary

This file combines two problematic patterns: a repository struct stores a `context.Context`, and the only public method performs a database read with `context.Background()` while never using the stored context. That makes request cancellation and deadlines invisible to the SQL call and misleads readers into thinking the constructor-bound context governs operations.

---

## Finding 1

【风险等级】 High

【问题代码行及解释】
`/Users/macmini/.claude/skills/go-context-audit/evals/files/struct_storage.go:8-17`

`UserRepository` holds `ctx context.Context` as a field, and `NewUserRepository` copies the caller’s context into that field. Long-lived or shared repository instances can carry a **stale or already-canceled** context, or blur **who owns cancellation** for future methods. Per the audit skill and `patterns.md`, struct-stored request context hides request lifetime in object state and encourages incorrect reuse across calls.

【修复建议代码】
```go
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}
```

Pass `ctx` on each method that does I/O (see Finding 2).

---

## Finding 2

【风险等级】 Medium

【问题代码行及解释】
`/Users/macmini/.claude/skills/go-context-audit/evals/files/struct_storage.go:20-22`

`LoadUser` does not accept a `context.Context`. It calls `QueryRowContext` with `context.Background()`, so the query **ignores** upstream cancellation, deadlines, and tracing propagation from the HTTP/RPC layer. The struct field `ctx` is **not used** here, which **hides** the mismatch: constructing `NewUserRepository(ctx, db)` suggests request scoping, but the read path is fully detached. Under load or client disconnect, work can continue longer than intended.

【修复建议代码】
```go
func (r *UserRepository) LoadUser(ctx context.Context, id string) error {
	return r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE id = $1", id).Err()
}
```

(After removing the struct field from Finding 1, `ctx` is always explicit at the call site.)

---

## Finding 3

【风险等级】 Low

【问题代码行及解释】
`/Users/macmini/.claude/skills/go-context-audit/evals/files/struct_storage.go:13-18`

`NewUserRepository` takes `ctx` as the first parameter while `db` is second. For constructors, putting `context.Context` first is acceptable; the larger issue is **storing** that `ctx`. If the constructor is kept for any reason, document that it must not be request-scoped—or remove it entirely in favor of per-method `ctx` (preferred).

【修复建议代码】
Prefer removing the constructor `ctx` parameter entirely (see Finding 1) so there is no ambiguous “default” context on the type.

---

## Residual review gaps

- Call sites of `NewUserRepository` and `LoadUser` were not reviewed; repository lifetime (per-request vs singleton) cannot be proven from this file alone.
- No goroutines, `WithCancel`/`WithTimeout`, or `WithValue` appear in this file; those dimensions are N/A here.

---

## No issues found (N/A categories)

- Goroutine lifecycle / `ctx.Done()` handling: none present.
- Derived contexts without `cancel`: none present.
- `context.WithValue` for business keys: none present.
