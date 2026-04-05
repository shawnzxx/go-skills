# Uncertainties / assumptions

- **Intended API contract for `jobs`:** Unclear whether producers are expected to close `jobs` to signal completion. The review flags the closed-channel behavior as a risk; if the contract is “never close,” severity of that item is lower in practice.
- **`handle` / `flush` behavior:** Bodies are empty stubs in this fixture; real implementations might spawn more goroutines or block. This review only covers what is visible in `goroutine_leak.go`.
- **Whether `Start` may be called multiple times:** Not specified; multiple calls would spawn multiple immortal goroutines unless guarded by the caller.
