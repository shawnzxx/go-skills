# Go concurrency review: `goroutine_leak.go`

**Task intent:** Review for hidden goroutine leaks, shutdown bugs, and context misuse. Findings below use **risk** (High / Medium / Low), **locations**, and **repair suggestions**.

---

## 1. Goroutine never terminates (leak) — **High**

**Location:** `Start` — goroutine started at lines 11–23; loop has no exit path.

**Issue:** The background goroutine runs an infinite `for` with `select` on `jobs` and `ticker.C` only. There is no `case <-ctx.Done():` (or any other shutdown signal), so the goroutine **outlives** cancellation of the passed-in `context.Context` and continues until the process exits.

**Repair suggestions:**

- Add `case <-ctx.Done(): return` (or `break` with labeled loop) so cancellation stops the goroutine.
- Optionally `defer ticker.Stop()` is already correct; keep it so the ticker is released when the goroutine exits.
- Document whether `Start` owns the goroutine lifecycle vs. caller expectations (e.g. “call `Stop` or cancel context to release”).

---

## 2. `context.Context` parameter unused (context misuse) — **High**

**Location:** `Start(ctx context.Context, jobs <-chan string)` — `ctx` is never referenced inside `Start` or the goroutine (lines 10–23).

**Issue:** Callers will reasonably assume `ctx` controls shutdown or deadlines for work started here. Ignoring it is **explicit context misuse**: cancellation does not propagate, and any timeout/deadline on `ctx` has no effect on this loop.

**Repair suggestions:**

- Use `ctx` in the `select` as above, and/or pass `ctx` into `handle` / `flush` if those perform I/O or should honor deadlines.
- If shutdown is intentionally not context-driven, remove `ctx` from the API and use a dedicated `Stop()` / `Close()` to avoid a misleading signature (prefer API honesty over unused parameters).

---

## 3. Closed `jobs` channel can cause a tight receive loop — **Medium**

**Location:** `case job := <-jobs:` (lines 17–18).

**Issue:** In Go, receiving from a **closed** channel yields the zero value **immediately** and repeatedly. With `select`, that case is always ready, so the loop can spend most of its time calling `p.handle("")` in quick succession, starving the ticker and burning CPU. This is a common shutdown/contract bug when producers close `jobs` to signal “no more work.”

**Repair suggestions:**

- Use the two-value form: `job, ok := <-jobs`; if `!ok`, `return` (or exit after flush), after defining desired semantics when the channel closes.
- Alternatively document that `jobs` must **not** be closed while this goroutine runs, and only signal completion via `ctx` or another channel.

---

## 4. No synchronization with caller after `Start` returns — **Low**

**Location:** `Start` returns immediately after `go func() { ... }()` (line 23).

**Issue:** This is normal for “fire background worker” APIs, but combined with missing shutdown, callers have no hook to `Wait()` for drain or exit unless they add one.

**Repair suggestions:**

- If needed, return a `WaitGroup`, `done` channel, or `Stop` func that closes a shutdown channel and waits for the goroutine to finish.

---

## Summary

| Risk   | Topic                          | Lines (approx.) |
|--------|--------------------------------|-----------------|
| High   | Goroutine leak, no shutdown  | 11–23           |
| High   | `ctx` ignored (misuse)       | 10, 11–23       |
| Medium | Closed `jobs` → busy loop     | 17–18           |
| Low    | No join/wait API              | 10–24           |

**Primary fix:** Combine `ctx.Done()` handling with a clear rule for `jobs` (close + exit, or never close until worker stopped).
