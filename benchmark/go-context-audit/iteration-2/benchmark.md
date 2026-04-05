# Benchmark Summary

- Skill: `go-context-audit`
- Iteration: `2`
- Runs per configuration: `1`

## Result

- `with_skill` pass rate mean: `1.00`
- `without_skill` pass rate mean: `0.83`
- Delta: `+0.17`

## Observations

- Iteration-2 is rerun from fresh real LLM outputs for all 6 runs.
- `goroutine-cancel-leak` strict `High` wording check causes baseline to miss one assertion (`Critical` wording mismatch).
- `withvalue-business-data` is less discriminative in this run because both configs produced strong fixes.
- Time and token metrics were not captured; fields remain placeholders.
