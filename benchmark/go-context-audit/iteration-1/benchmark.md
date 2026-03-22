# Benchmark Summary

- Skill: `go-context-audit`
- Iteration: `1`
- Runs per configuration: `1`

## Result

- `with_skill` pass rate mean: `1.00`
- `without_skill` pass rate mean: `0.83`
- Delta: `+0.17`

## Observations

- The skill produced stronger outputs on the two most discriminating evals.
- In `goroutine-cancel-leak`, the skill reliably produced a concrete `ctx.Done()` repair snippet while baseline only described the fix in prose.
- In `withvalue-business-data`, the skill more clearly classified `userID` and `orderID` as business parameters that should not be hidden in context, while baseline treated the pattern more ambiguously.
- `struct-storage-and-background` did not differentiate the configurations in this round because both baseline and with-skill outputs were already strong.
- Time and token metrics were not captured for this iteration, so benchmark timing fields are placeholders only.
