# Context audit: `withvalue_business.go`

**Task intent:** Audit this Go snippet for context misuse. Assess whether `context.WithValue` is used appropriately; include severity, explanation, and corrected code.

**Scope:** Single file `evals/files/withvalue_business.go` only (no transitive callees inspected).

---

## Executive summary

`userID` and `orderID` are required business inputs stored and read via `context.WithValue` / `ctx.Value`. That pattern is **not appropriate** for this data: it hides dependencies, complicates testing, and encourages panics from unchecked type assertions. `traceID` as cross-cutting metadata in context is **generally acceptable**, but mixing it with business IDs in one helper blurs boundaries.

---

## Finding 1

【风险等级】 Medium

【问题代码行及解释】

`withvalue_business.go:13-18`, `withvalue_business.go:20-23`

`BuildContext` attaches `userID` and `orderID` to the context with `context.WithValue`. `PlaceOrder` then pulls them out with `ctx.Value(...).(string)`. Per Go community guidance and this audit’s risk model, identifiers and entities required for business logic (user ID, order ID) should be **explicit function parameters** (or a typed request struct), not context values.

Failure modes and drawbacks:

- Callers cannot see that `PlaceOrder` depends on those values without reading its body; static analysis and refactoring suffer.
- `ctx.Value(userIDKey).(string)` panics if the key is missing or the stored type is not `string`—easy to misuse when contexts are built in different packages or tests.
- Tests must reconstruct a “magic” context instead of calling APIs with plain values.

`traceID` in the same builder is closer to acceptable cross-cutting metadata (see Finding 2), but bundling it with business IDs encourages treating context as an implicit parameter bag.

【修复建议代码】

Keep trace/correlation in context if you want; pass business IDs explicitly:

```go
package auditfixture

import "context"

type contextKey string

const traceIDKey contextKey = "traceID"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func PlaceOrder(ctx context.Context, userID, orderID string) error {
	traceID, _ := ctx.Value(traceIDKey).(string)
	_, _, _ = userID, orderID, traceID
	return nil
}
```

If you prefer zero context values for production code, pass `traceID` as a field on a small options struct or logger instead of `WithValue`.

---

## Finding 2

【风险等级】 Low

【问题代码行及解释】

`withvalue_business.go:21-23`

`ctx.Value(...).(string)` uses a **single-result type assertion**. If `BuildContext` was not called, keys diverged, or a different type was stored, the program panics. This is not a goroutine/timer leak, but it is a real robustness issue tied to `WithValue` usage.

【修复建议代码】

Use the comma-ok form and handle absence explicitly (return error or use defaults deliberately):

```go
func PlaceOrder(ctx context.Context, userID, orderID string) error {
	traceID, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		traceID = "" // or return fmt.Errorf("missing trace ID in context")
	}
	_, _, _ = userID, orderID, traceID
	return nil
}
```

(Prefer moving `userID`/`orderID` out of context as in Finding 1; then missing-trace handling is the main `Value` concern.)

---

## Finding 3 (informational — not a misuse)

【风险等级】 (none — acceptable pattern)

【问题代码行及解释】

`withvalue_business.go:16`

Storing `traceID` via `context.WithValue` is **often appropriate** when it is used for logging, tracing, or correlation across middleware and handlers. Unexported typed keys (`type contextKey string`) are a good practice to avoid collisions.

No change required for trace ID **in isolation**; the main issue is pairing it with business IDs in the same API surface (Finding 1).

---

## Residual gaps

- Call sites of `BuildContext` / `PlaceOrder` were not reviewed; ownership of cancellation and whether contexts cross package boundaries could not be fully verified from this file alone.
