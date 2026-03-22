# go-skills

Agent Skills for Go code review, mistake prevention, and engineering heuristics.

This repository turns real Go review lessons and production mistakes into reusable Agent Skills. The goal is to help coding agents catch high-signal issues such as context misuse, lifecycle leaks, hidden coupling, and API design mistakes before they reach production.

## Why this repo exists

Many Go failures are not syntax errors. They are ownership bugs, cancellation leaks, shutdown mistakes, and implicit business dependencies that often survive tests and only surface during code review or in production.

This repo packages those lessons as portable skills so they can be reused across projects and shared with other Go developers.

## Repository structure

```text
go-skills/
├── docs/
├── skills/
│   └── go-context-audit/
└── README.md
```

## Installation

### `npx skills`

```sh
npx skills add git@github.com:shawnzxx/go-skills.git
```

### Manual

Clone this repository and copy the `skills/` directory into your agent's skill path, or place the repository in a location your agent can discover.

For Claude Code, one common option is to copy a skill directory into `~/.claude/skills/`.

## Skills

| Skill              | Description                                                                                                             |
| ------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| `go-context-audit` | Reviews Go code for context leaks, cancellation mistakes, goroutine lifetime issues, and misuse of `context.WithValue`. |

## Publishing path

This repository is designed to follow the standard Agent Skills format first, then support marketplace-style distribution on top of that foundation.

That means the repo should always remain:

- easy to clone and inspect
- easy to install manually
- easy to validate as a standard skill collection
- easy to register with marketplaces or plugin registries later

## Roadmap

Near-term candidates:

- `go-error-handling-review`
- `go-channel-lifecycle-audit`
- `go-timeout-boundary-review`
- `go-interface-design-review`
- `go-db-transaction-audit`
- `go-test-smell-review`

## Contributing

Contributions are welcome, especially when they add:

- narrowly scoped review skills
- concrete false-positive boundaries
- realistic eval fixtures based on actual mistakes
- examples that improve trigger accuracy

## License

MIT
