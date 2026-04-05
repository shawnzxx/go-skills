# Uncertainties / assumptions

1. **Auth / `userID`:** Whether storing `userID` in `context.Context` is acceptable depends on project conventions (middleware-only vs explicit parameters). Marked as low/situational in the review.

2. **Empty strings:** The sample “corrected” code treats empty `userID` / `traceID` as invalid; real systems may allow empty trace IDs or use different validation.

3. **OpenTelemetry:** Suggestion to align with OTel spans is optional; not assumed to be in use for this fixture.

4. **Scope:** Only the provided file was reviewed; no callers of `BuildContext` / `PlaceOrder` were inspected, so whether context is always built via `BuildContext` is unknown — this strengthens the case against panicking assertions.
