# Uncertainties (single-file review)

1. **Repository lifetime:** It is unknown whether `*UserRepository` is created per HTTP request, per connection, or as a process singleton. If it were strictly per-request and always discarded before the request context completes, struct-stored `ctx` would still be non-idiomatic and fragile, but the “stale context on reuse” failure mode would be less likely in practice.

2. **Intended semantics of `NewUserRepository(ctx, …)`:** The stored `ctx` is never read in the visible code. It is unclear whether future methods were planned to use it (which would worsen coupling) or whether it is dead API surface—either way, current behavior does not match a reader’s likely expectation.

3. **Transitive behavior:** Without caller/callee files, we cannot confirm whether any parent already wraps the DB with timeouts or whether tracing middleware depends on context propagation into `database/sql`.
