# Context misuse review: `withvalue_business.go`

**Task intent:** Audit for `context.WithValue` appropriateness; provide severity, explanation, and corrected code.

**Scope:** `auditfixture` package — `BuildContext`, `PlaceOrder`.

---

## Summary

`context.WithValue` is used for `userID`, `orderID`, and `traceID`. The **typed, unexported key type** (`type contextKey string`) is appropriate and avoids key collisions across packages. **`traceID`** is a reasonable request-scoped, cross-cutting value. **`userID`** is a common (though debatable) pattern when threaded through middleware. **`orderID`** is a poor fit for context when it is really **per-operation data** for `PlaceOrder` — the [Go blog on context](https://go.dev/blog/context) recommends using context values only for data that **transits processes and APIs**, not as a substitute for function parameters.

The **highest-severity issue** is **unchecked type assertions** on `ctx.Value`, which **panic** if the key is missing or the value is the wrong type.

---

## Findings

### 1. Unchecked `ctx.Value` type assertions — **High**

| Item | Detail |
|------|--------|
| **Severity** | High |
| **Location** | `PlaceOrder`: `ctx.Value(...).(string)` for `userID`, `orderID`, `traceID` |
| **Explanation** | A single missing key or non-`string` value causes a **runtime panic**. Context values are not type-safe; callers can pass a context built without `BuildContext`, tests can omit values, or future refactors can break silently. |
| **Recommendation** | Use the comma-ok form, return a structured error, or validate once at the boundary (e.g. middleware) and pass typed data explicitly downstream. |

### 2. `orderID` carried only via context — **Medium**

| Item | Detail |
|------|--------|
| **Severity** | Medium |
| **Location** | `BuildContext` / `PlaceOrder` |
| **Explanation** | `orderID` is domain data for a specific operation. Hiding it inside `context.Context` obscures the API of `PlaceOrder`, makes call sites harder to read, and encourages “stringly typed” context bags. Prefer `PlaceOrder(ctx, orderID string)` (and optionally other explicit parameters) when the value is not truly cross-cutting. |
| **Recommendation** | Pass `orderID` as a function argument; keep context for cancellation, deadlines, and values that many layers need without changing every signature (e.g. trace ID). |

### 3. `userID` in context — **Low / situational**

| Item | Detail |
|------|--------|
| **Severity** | Low (or Medium if your style guide forbids auth identity in context) |
| **Location** | `BuildContext` / `PlaceOrder` |
| **Explanation** | Request-scoped principal identity is often stored in context after authentication middleware. That can be justified, but alternatives include explicit structs or parameter objects for handlers/services. Severity depends on team conventions. |
| **Recommendation** | If you keep it, document the contract, use typed accessors, and never panic on read (see finding 1). |

### 4. `traceID` in context — **Informational (appropriate use)**

| Item | Detail |
|------|--------|
| **Severity** | Informational |
| **Location** | `BuildContext` / `PlaceOrder` |
| **Explanation** | Trace/correlation IDs are classic request-scoped values that propagate through many layers without threading a trace parameter everywhere. |
| **Recommendation** | Consider aligning with OpenTelemetry (`trace.SpanFromContext`) if you adopt full tracing; otherwise this pattern is reasonable. |

### 5. Chained `WithValue` — **Low**

| Item | Detail |
|------|--------|
| **Severity** | Low |
| **Location** | `BuildContext` |
| **Explanation** | Multiple `WithValue` calls add a small chain; acceptable for a few keys. Very large context bags hurt clarity and performance slightly. |
| **Recommendation** | If many keys accumulate, consider a small request-scoped struct stored once under a single private key (still with safe extraction). |

---

## Corrected code (illustrative)

Goals: (1) no panic on missing/wrong types, (2) `orderID` as an explicit parameter, (3) keep trace (and optionally user) in context with safe reads.

```go
package auditfixture

import (
	"context"
	"errors"
	"fmt"
)

type contextKey string

const (
	userIDKey  contextKey = "userID"
	traceIDKey contextKey = "traceID"
)

func BuildContext(ctx context.Context, userID, traceID string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, userID)
	ctx = context.WithValue(ctx, traceIDKey, traceID)
	return ctx
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok && v != ""
}

func TraceIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(traceIDKey).(string)
	return v, ok && v != ""
}

func PlaceOrder(ctx context.Context, orderID string) error {
	userID, ok := UserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("auditfixture: missing userID in context: %w", errors.New("invalid context"))
	}
	traceID, ok := TraceIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("auditfixture: missing traceID in context: %w", errors.New("invalid context"))
	}

	_, _, _ = userID, orderID, traceID
	return nil
}
```

**Note:** Error messages and whether empty strings are invalid are policy choices; adjust to match your API.

---

## Verdict on `context.WithValue` “appropriateness”

- **Appropriate:** `traceID` (and often `userID` in middleware-heavy stacks), with **safe retrieval**.
- **Questionable / usually inappropriate as the only carrier:** `orderID` for `PlaceOrder` — prefer an explicit parameter unless the same context is used uniformly across a pipeline where `orderID` is truly request-global (rare).

---

*Review produced without using the `go-context-audit` skill file, per evaluation constraints.*
