---
name: go-context-audit
description: Review Go code for context leaks, goroutine lifecycle issues, cancellation mistakes, and suspicious context propagation. Use when auditing `context.Context`, `ctx` propagation, timeout handling, goroutine leaks, or `context.WithValue` misuse in Go code.
license: MIT
metadata:
  repository: go-skills
  category: go-review
---

# Go Context Audit

Review Go code for `context.Context` leaks and misuse, then report findings with concrete repair guidance.

## When To Use

Use this skill when the user asks to:

- review Go code for context issues
- check context leaks or cancellation bugs
- audit `ctx` propagation or timeout handling
- review goroutine or concurrency issues in Go where context misuse may be involved

Do not use this skill for non-Go concurrency reviews or for requests that only need general style feedback.

## Review Scope

Start with the smallest relevant scope:

1. User-specified file or snippet
2. Current diff or files under review
3. Directly connected callers and callees
4. Wider repository scan only if the user asks for repo-wide auditing or the local evidence suggests a broader pattern

This keeps the review focused and prevents noisy repo-wide findings when the user only wants feedback on a narrow change.

## What To Look For

### High-signal findings

- `context.Context` stored in structs, long-lived objects, caches, or worker state
- goroutines that may outlive the request or parent operation without checking `ctx.Done()`
- loops, selects, or blocking calls that have no cancellation path
- `context.WithCancel`, `context.WithTimeout`, or `context.WithDeadline` created without a matching `cancel`
- code that drops the upstream `ctx` and replaces it with `context.Background()` or `context.TODO()`
- `context.WithValue` used for required business inputs such as user IDs, tenant IDs, order IDs, auth state, pagination, or feature flags

### Supporting checks

- `ctx` is not the first parameter
- downstream calls do not receive the upstream `ctx`
- background workers launch from request handlers without a clear ownership boundary
- cleanup paths exist but are not guaranteed to run on error or cancellation
- channel-driven goroutines rely on channel close for shutdown but do not check `value, ok := <-ch` or otherwise make exit semantics explicit

## Risk Model

### High

Use `High` when the code can leak goroutines, timers, or request-scoped work, or when it stores request context in long-lived state.

Examples:

- a goroutine loops forever and never listens to `ctx.Done()`
- a `WithTimeout` context is created inside a hot path and `cancel` is never called
- a service or client struct stores a request-derived `ctx`

### Medium

Use `Medium` when the code obscures cancellation or smuggles core business data through context, but the immediate leak risk is indirect.

Examples:

- `context.WithValue` carries `userID`, `tenantID`, or `orderID`
- a helper replaces the incoming `ctx` with `context.Background()`
- a downstream dependency call silently omits the upstream `ctx`

### Low

Use `Low` for convention issues or clarity problems that reduce maintainability but are unlikely to leak by themselves.

Examples:

- `ctx` is not the first parameter
- cancellation is technically correct but harder to follow than necessary

## False Positive Boundaries

Be careful not to over-report:

- `context.WithValue` is often acceptable for tracing IDs, request IDs, log correlation IDs, span state, and other cross-cutting metadata
- auth or identity data may appear in context in middleware-heavy stacks, but if values such as `userID` or `tenantID` are required business inputs for the reviewed function, default to flagging them unless the cross-cutting ownership is explicit and well-documented
- long-lived background services may intentionally derive a fresh root context at process startup; flag this only when request-scoped work is being detached without an explicit ownership decision
- some goroutines exit because their input channel closes; if there is a clear, bounded owner and shutdown path, explain the tradeoff instead of assuming a leak

If uncertain, state the uncertainty explicitly and explain what evidence would confirm or refute the finding.

For more examples of boundaries and detection patterns, read `references/patterns.md`.

## Review Workflow

1. Confirm the review scope from the user's files, diff, or snippet.
2. Trace where each relevant `ctx` comes from and where it should end.
3. Inspect goroutine launch sites, loops, blocking operations, and timeout helpers.
4. Check whether derived contexts are canceled on every exit path.
5. Inspect `WithValue` call sites and judge whether the data is cross-cutting metadata or required business input.
6. Rank findings by risk and provide the smallest practical fix that restores clear ownership and cancellation behavior.

## Output Format

Always return Markdown. For each finding, use this structure:

````md
## Finding N

【Risk】 High

【Code And Explanation】
`path/to/file.go:12-24`

Explain what the code is doing, why it is risky, and what failure mode it can trigger.

【Suggested Fix】

```go
// minimal fix here
```
````

If there are no findings, say so explicitly and mention any residual review gaps, such as "did not inspect transitive callees" or "could not verify shutdown ownership from the provided snippet".
Do not replace `【Suggested Fix】` with prose only when a practical fix can be shown. Prefer a real Go code block, even if it is a minimal patch sketch.

## Fix Guidance

Prefer minimal, concrete fixes:

- move `ctx` from struct fields to method parameters
- add `select { case <-ctx.Done(): return }` or equivalent cancellation handling around long-running goroutines
- if shutdown may happen by closing a channel, show `value, ok := <-ch` and an explicit `if !ok { return }` branch when that is part of the intended lifecycle
- `defer cancel()` immediately after deriving a cancelable context, unless ownership is intentionally transferred and documented
- replace `WithValue` business data with explicit parameters, typed request objects, or dedicated option structs
- propagate the incoming `ctx` instead of creating a fresh root context in library code

When showing code, keep the example tightly scoped to the finding rather than rewriting the whole function.

## Example Findings

### Example 1: Context stored in a struct

```go
type Service struct {
    ctx context.Context
    db  *sql.DB
}
```

This usually deserves at least `Medium`, and often `High` if the struct outlives a single request. Explain that request-scoped cancellation is now hidden inside object state and can be reused incorrectly across calls.

### Example 2: Goroutine without cancellation

```go
go func() {
    for msg := range jobs {
        process(msg)
    }
}()
```

If the goroutine is tied to request work or a derived context, flag the missing cancellation path and suggest a `select` that honors `ctx.Done()` or another explicit shutdown signal.

### Example 3: Business data in WithValue

```go
ctx = context.WithValue(ctx, userIDKey, userID)
```

Flag this when `userID` is required business input rather than cross-cutting metadata. Recommend passing it explicitly.

## Review Style

- Lead with the highest-risk finding.
- Explain why the pattern is dangerous, not just that it is non-idiomatic.
- Prefer evidence over speculation.
- Mention uncertainty when ownership or shutdown behavior cannot be proven from the visible code.
- Default to review output, not direct code edits, unless the user explicitly asks for a patch.
