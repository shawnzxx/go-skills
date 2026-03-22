# Detection Patterns

Use this reference when the main review needs more nuance about likely bugs versus acceptable tradeoffs.

## Struct-Stored Context

Usually report when:

- a struct field has type `context.Context`
- a constructor captures a request-derived `ctx` and stores it for later calls
- worker state, repositories, or clients keep `ctx` across method boundaries

Why it matters:

- request lifetime becomes hidden object state
- cancellation ownership is no longer obvious
- future calls may accidentally reuse a stale or canceled context

Safer alternatives:

- pass `ctx` as a method parameter
- store only immutable configuration in the struct
- if process-level shutdown is needed, store a clearly named root context that is created at startup and not derived from a request

## Goroutine Lifetime And Cancellation

Likely leak patterns:

- `go func()` starts request-scoped work with no `ctx.Done()` handling
- ticker or timer loops never stop on cancellation
- goroutine reads from channels forever without a documented shutdown owner
- fan-out work derives child contexts but does not cancel them

Signals that reduce risk:

- `select` includes `case <-ctx.Done():`
- owner closes the input channel and that ownership is obvious from the surrounding code
- `errgroup.Group` or another structured concurrency helper controls lifecycle
- process-level background worker has a documented root context and shutdown path

## Derived Contexts Without Cancel

Flag when code creates:

- `context.WithCancel`
- `context.WithTimeout`
- `context.WithDeadline`
- `context.WithCancelCause`
- `context.WithTimeoutCause`
- `context.WithDeadlineCause`

and the corresponding cancel function is never called or ownership is unclear.

Common fixes:

- `ctx, cancel := context.WithTimeout(...); defer cancel()`
- return both `ctx` and `cancel` only when the caller is clearly taking ownership
- avoid deriving a context if the callee can accept the original one directly

## WithValue Boundaries

Usually acceptable:

- trace IDs
- span context
- request IDs
- log correlation metadata
- auth or instrumentation metadata that is truly cross-cutting and optional to most business calls

Borderline case:

- authenticated principal data such as `userID` may be injected by middleware, but if the reviewed function cannot do its business work without that value, treat it as a business dependency by default and explain why

Usually problematic:

- `userID`
- `tenantID`
- `accountID`
- `orderID`
- pagination or filtering data
- role, permission, or feature-gate decisions required for business behavior

Reasoning:

- required business inputs should be explicit in function signatures
- hidden parameters make dependencies harder to see, test, and refactor

## Replacing Upstream Context

Suspicious patterns:

- helper functions call `context.Background()` instead of using the incoming `ctx`
- library code silently uses `context.TODO()` in a call chain that already had a context
- request handlers detach work to a fresh background context without documenting ownership transfer

Potentially acceptable:

- process startup code creating a root context for daemon-level background workers
- intentionally detached fire-and-forget work with a documented lifecycle and durability strategy

When in doubt, explain the ownership question instead of asserting a bug without evidence.

## Channel-Close Shutdown Semantics

When a goroutine reads from a channel in a loop, also inspect whether channel close is intended to stop the worker.

Risky pattern:

- `case item := <-jobs:` with no `ok` check

Why it matters:

- reads from a closed channel return the zero value immediately
- the goroutine may spin, keep calling handlers with zero values, or mask shutdown bugs

Typical fix:

- `case item, ok := <-jobs: if !ok { return }`
