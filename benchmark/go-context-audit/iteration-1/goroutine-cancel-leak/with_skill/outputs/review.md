# Go concurrency review: `goroutine_leak.go`

**Task intent:** Please review this Go code for concurrency issues. I am worried there may be hidden goroutine leaks or shutdown bugs. If context misuse is involved, call it out explicitly. Return Markdown findings with risk, code locations, and repair suggestions.

**Scope:** Single file `evals/files/goroutine_leak.go` (package `auditfixture`) as provided; callers and channel ownership are not visible in this snippet.

---

## Finding 1

【风险等级】 High

【问题代码行及解释】
`goroutine_leak.go:10-23`

`Start` accepts `ctx context.Context` but never uses it. A goroutine runs an unbounded `for { select { ... } }` that only waits on `jobs` and a 5s ticker. There is no `case <-ctx.Done(): return` (or equivalent), so parent cancellation or request shutdown does not stop this worker. That is a classic pattern for **goroutines that outlive the operation** that was supposed to bound them.

If `jobs` is closed, receive on a closed channel yields the zero value (`""`) immediately and repeatedly, so the loop can keep spinning between empty jobs and ticker ticks instead of exiting—**shutdown via closing `jobs` alone is not reliable** without an explicit break or `ctx` branch.

Failure modes: leaked goroutine after shutdown; continued `handle`/`flush` calls after the caller believes work is stopped; harder-to-reason resource lifetime.

【修复建议代码】
```go
func (p *Processor) Start(ctx context.Context, jobs <-chan string) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}
				p.handle(job)
			case <-ticker.C:
				p.flush()
			}
		}
	}()
}
```

(Adjust `ok`/close semantics if the design intentionally keeps `jobs` open for the process lifetime; still keep `ctx.Done()` for bounded shutdown.)

---

## Finding 2

【风险等级】 Medium

【问题代码行及解释】
`goroutine_leak.go:10`

The API surface suggests cancellation-aware startup (`ctx` is the first parameter after the receiver), but the implementation **drops the upstream context** entirely. Downstream work (`handle`, `flush`) cannot be tied to the same cancellation/deadline tree without further propagation. Even if those methods later accept `ctx`, this goroutine would need to pass `ctx` (or a derived timeout) into them.

【修复建议代码】
```go
func (p *Processor) handle(ctx context.Context, job string) {}
func (p *Processor) flush(ctx context.Context) {}

// Inside the loop, after select receives work or ticker:
// p.handle(ctx, job) and p.flush(ctx)
```

(Exact signatures depend on whether `handle`/`flush` should honor per-job timeouts; at minimum, document why `ctx` is ignored if it is intentional.)

---

## Summary

| Finding | Risk   | Theme                                      |
|---------|--------|--------------------------------------------|
| 1       | High   | Unbounded goroutine; ignored `ctx.Done()`; fragile shutdown on channel close |
| 2       | Medium | Request/parent `ctx` unused; no propagation to work |

No `WithCancel`/`WithTimeout` leak of a `cancel` function appears in this file. No `WithValue` or struct-stored `context.Context` in this snippet.

**Residual gaps:** Caller behavior (who closes `jobs`, when `Start` is called, process vs request scope) is not in-file; confirming “leak in production” would require those call sites.
