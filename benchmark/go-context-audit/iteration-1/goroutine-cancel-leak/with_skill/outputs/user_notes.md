# Uncertainties (out-of-file evidence)

- **Callers of `Start`:** Whether `jobs` is closed on shutdown, whether `Start` is invoked once per process or per request, and whether a `WaitGroup` or other join mechanism exists elsewhere cannot be determined from `goroutine_leak.go` alone. These affect how severe a “leak” is in deployment (e.g. one leaked goroutine per process vs per HTTP request).
- **`handle` / `flush`:** Bodies are empty stubs; unknown if real implementations block, spawn more goroutines, or need their own deadlines—recommendations assume they should respect the same cancellation as the worker loop once implemented.
- **Intent:** If the design is a deliberate long-lived daemon with a root context only used by outer callers, the **High** finding still applies to *this* snippet because the passed-in `ctx` is unused; the fix might be renaming/documenting a process-level context vs request context, which would need product/architecture context not present in the file.
