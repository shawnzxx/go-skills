---
name: go-context-audit
description: Review Go code for context bugs, goroutine leaks, and memory issues. Use when user mentions Go code problems with context, goroutines, background workers, or asks to audit/debug code involving ctx. Catches leaks, missing cancel(), stored ctx, time.After in loops, and WithValue anti-patterns.
license: MIT
metadata:
  repository: go-skills
  category: go-review
---

# Go Context Audit

Review Go code for `context.Context` misuse and report findings with risk-ranked, actionable repair guidance.

Context bugs are insidious because they compile, pass basic tests, and only surface under production load as goroutine leaks, stale data, or cascading timeouts. This skill helps catch them early.

## When To Use

Use this skill when the user asks to:

- review Go code for context issues, leaks, or cancellation bugs
- audit `ctx` propagation, timeout handling, or deadline management
- check goroutine or concurrency issues in Go (context misuse is often the root cause)
- investigate memory growth, goroutine count increases, or resource leaks in Go services
- review middleware, HTTP handler, or gRPC interceptor context patterns

Do not use for non-Go concurrency reviews or purely stylistic Go feedback.

## Review Scope

Start with the smallest relevant scope, then expand only when evidence warrants it:

1. User-specified file or snippet
2. Current diff or files under review
3. Directly connected callers and callees (one hop)
4. Wider repository scan only when the user requests it or local findings suggest a systemic pattern

## What To Look For

### Critical Patterns (almost always bugs)

- `context.Context` stored in struct fields, caches, or long-lived objects
- goroutines that outlive their parent operation without checking `ctx.Done()`
- `context.WithCancel` / `WithTimeout` / `WithDeadline` created without calling `cancel`
- code that drops upstream `ctx` and replaces it with `context.Background()` or `context.TODO()` inside library or handler code
- `time.After()` inside a `select` loop — each iteration allocates a timer that cannot be garbage collected until it fires, causing memory leaks under load
- `context.WithValue` used for required business inputs (user IDs, tenant IDs, order IDs, auth decisions, pagination)

### Secondary Patterns (often problematic, context-dependent)

- `ctx` is not the first parameter of a function
- downstream calls do not receive the upstream `ctx`
- background workers launched from request handlers without a clear ownership boundary
- `errgroup.WithContext` used but the original (non-derived) context is passed to child goroutines, bypassing the group's cancellation
- cleanup paths that are not guaranteed to run on error or cancellation
- `select` with a `default` case that creates a hot-spin loop when channels are empty
- channel-driven goroutines that rely on channel close but lack `value, ok := <-ch` for safe exit
- `context.WithoutCancel` (Go 1.21+) used to detach work — intentional but deserves scrutiny for lifetime ownership

## Risk Model

### High

The code can leak goroutines, timers, or request-scoped resources, or stores request context in long-lived state.

Examples:
- a goroutine loops forever without listening to `ctx.Done()`
- `WithTimeout` is created in a hot path and `cancel` is never called
- a service struct stores a request-derived `ctx`
- `time.After` is used inside a `for`/`select` loop (timer leak on every iteration)

### Medium

The code obscures cancellation, smuggles business data through context, or silently drops the upstream context, but the immediate leak risk is indirect.

Examples:
- `context.WithValue` carries `userID`, `tenantID`, or `orderID`
- a helper replaces incoming `ctx` with `context.Background()`
- `errgroup.WithContext` is used but children receive the wrong ctx
- `context.WithoutCancel` detaches work without documented ownership

### Low

Convention issues that reduce clarity but are unlikely to leak by themselves.

Examples:
- `ctx` is not the first parameter
- cancellation is technically correct but harder to follow than necessary
- `context.TODO()` appears in code that should have a real context by now

## False Positive Boundaries

Avoid over-reporting:

- `context.WithValue` is acceptable for trace IDs, request IDs, log correlation, and span context
- auth/identity data in context is common in middleware stacks — flag only when the reviewed function requires it as a business input and cannot operate without it
- process-level background services that derive a fresh root context at startup are fine; flag only when request-scoped work is being detached without explicit ownership
- goroutines that exit because their input channel closes are OK if there is a clear, bounded owner and shutdown path — explain the tradeoff rather than assuming a leak

When uncertain, state the uncertainty and what evidence would confirm or refute the finding.

For detailed detection heuristics and edge cases, read `references/patterns.md`.

## Review Workflow

1. **Scope**: confirm the review target from the user's files, diff, or snippet.
2. **Trace origins**: for each `ctx`, trace where it was created and where it should be canceled.
3. **Inspect launch sites**: goroutine spawns, loops, blocking ops, timer/ticker usage.
4. **Check cancellation paths**: every derived context must have its `cancel` called on all exit paths.
5. **Audit `WithValue`**: distinguish cross-cutting metadata from business-required inputs.
6. **Assess lifetime boundaries**: do goroutines, timers, or deferred work outlive their parent context?
7. **Rank and report**: order findings by risk, provide minimal fixes.

## Output Format

Return Markdown. For each finding:

````md
## Finding N

**Risk:** High | Medium | Low

**Location:** `path/to/file.go:12-24`

**Issue:**
Explain what the code does, why it is risky, and what failure mode it can trigger under load or cancellation.

**Fix:**

```go
// minimal targeted fix
```
````

If there are no findings, say so explicitly and note any review gaps (e.g., "did not inspect transitive callees" or "could not verify shutdown ownership from the provided snippet").

Always prefer showing a real Go code fix over prose-only guidance.

## Fix Guidance

Prefer minimal, concrete fixes:

- move `ctx` from struct fields to method parameters
- add `select { case <-ctx.Done(): return }` around long-running goroutine loops
- replace `time.After()` in loops with a reusable `time.NewTimer` + `timer.Reset` pattern
- `defer cancel()` immediately after deriving a cancelable context (unless ownership is intentionally transferred and documented)
- replace `WithValue` business data with explicit parameters, typed request objects, or option structs
- propagate the incoming `ctx` instead of creating a fresh root context in library code
- use `value, ok := <-ch` and return on `!ok` for channel-close shutdown semantics
- when using `errgroup.WithContext`, pass the derived context (not the parent) to child goroutines

Keep fix examples tightly scoped to the finding — do not rewrite the entire function.

## Examples

### Example 1: Context stored in a struct

```go
type Service struct {
    ctx context.Context
    db  *sql.DB
}
```

At least Medium risk, often High if the struct outlives a single request. The request's cancellation becomes hidden object state and can be reused incorrectly across calls. Fix: pass `ctx` as a method parameter.

### Example 2: Goroutine without cancellation

```go
go func() {
    for msg := range jobs {
        process(msg)
    }
}()
```

If the goroutine is tied to request work or a derived context, flag the missing cancellation path. Fix: add a `select` that honors `ctx.Done()`.

### Example 3: Business data in WithValue

```go
ctx = context.WithValue(ctx, userIDKey, userID)
```

Flag when `userID` is a required business input. Cross-cutting metadata like trace IDs is acceptable. Fix: pass `userID` as an explicit function parameter.

### Example 4: Timer leak in select loop

```go
for {
    select {
    case msg := <-ch:
        handle(msg)
    case <-time.After(5 * time.Second):
        flush()
    }
}
```

Each loop iteration allocates a new timer that cannot be GC'd until it fires. Under high throughput this leaks memory. Fix: use `time.NewTicker` or `time.NewTimer` with `Reset`.

## Review Style

- Lead with the highest-risk finding.
- Explain **why** the pattern is dangerous, not just that it is non-idiomatic.
- Prefer evidence over speculation.
- Mention uncertainty when ownership or shutdown behavior cannot be proven from visible code.
- Default to review output, not direct code edits, unless the user explicitly asks for a patch.
