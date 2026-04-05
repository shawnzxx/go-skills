# Execution log — go-context-audit eval (withvalue_business)

## Steps

1. Read `/Users/macmini/.claude/skills/go-context-audit/SKILL.md` for review scope, risk model, output format, and `WithValue` guidance.
2. Read `/Users/macmini/.claude/skills/go-context-audit/references/patterns.md` for `WithValue` boundaries (metadata vs business data) and false-positive notes.
3. Read target file `evals/files/withvalue_business.go` and traced context flow: `BuildContext` derives a context with three `WithValue` keys; `PlaceOrder` reads all three via `ctx.Value` and string type assertions.
4. Classified keys: `userID` / `orderID` → required business inputs (problematic in context per skill); `traceID` → cross-cutting metadata (generally OK).
5. Checked for other high-signal issues from the skill (struct-stored context, goroutines, `WithCancel`/`WithTimeout` without cancel, `Background` replacement): **none present** in this file.
6. Wrote `review.md` with severity-tagged findings, explanations, and minimal corrected Go snippets; logged uncertainties in `user_notes.md`.

## Conclusion

**`context.WithValue` is not used appropriately for `userID` and `orderID`** — treat as **Medium** severity (smuggled business parameters). **`traceID` in context is acceptable** in isolation. **Low** severity: unchecked type assertions on `ctx.Value` can panic. No High-severity leak patterns in the provided snippet.
