# Roadmap

## Current scope

- Publish reusable Agent Skills for Go code review
- Capture real mistakes as portable review heuristics
- Keep each skill narrow, composable, and easy to trigger correctly

## Current skills

- `go-context-audit`

## Next skill candidates

- `go-error-handling-review`
- `go-channel-lifecycle-audit`
- `go-timeout-boundary-review`
- `go-interface-design-review`
- `go-db-transaction-audit`
- `go-test-smell-review`

## Skill design rules

- One skill should solve one review problem well.
- Skill descriptions should explain both what the skill does and when to use it.
- Prefer concrete review signals over generic style advice.
- Keep `SKILL.md` concise and move deeper material into `references/`.
- Add eval fixtures whenever a pattern is subtle or easy to regress.

## Publishing milestones

### Milestone 1

- Make the repository usable as a standard Agent Skills repo
- Support manual installation and Git-based installation

### Milestone 2

- Add at least 3 well-scoped Go review skills
- Improve eval coverage and examples

### Milestone 3

- Add marketplace-specific metadata or registry integration as needed
- Document release and compatibility expectations
