# Migrate a Local Skill

This guide explains how to move a local skill into this repository and prepare it for public sharing.

## Goal

Move a skill from a local agent path such as `~/.claude/skills/<skill-name>` into this repository under `skills/<skill-name>`.

## Step 1: Create the target directory

Use this layout:

```text
skills/<skill-name>/
├── SKILL.md
├── references/
├── evals/
└── scripts/
```

Only `SKILL.md` is required. The other directories are optional.

## Step 2: Copy the existing skill content

Move the current `SKILL.md` into:

```text
skills/<skill-name>/SKILL.md
```

If the skill already has helper files such as references, examples, evals, or scripts, keep the same relative structure when possible.

## Step 3: Verify frontmatter

Check these fields:

- `name`
- `description`

Recommended additions:

- `license`
- `metadata`

Rules:

- the `name` must match the parent directory name
- the `description` should describe both what the skill does and when to use it
- keep the wording specific so agents can trigger the skill correctly

## Step 4: Tighten the description

A good description should include:

- the target language or domain
- the main problem class
- trigger phrases the user is likely to say

Example:

```yaml
description: Review Go code for context leaks, cancellation mistakes, and goroutine lifecycle issues. Use when auditing `context.Context`, `ctx` propagation, timeouts, goroutine leaks, or `context.WithValue` misuse in Go code.
```

## Step 5: Keep `SKILL.md` concise

Prefer:

- activation conditions
- review workflow
- output format
- a few concrete examples

Move deeper nuance into `references/` when the file becomes too long.

## Step 6: Add eval fixtures

If the skill captures subtle review heuristics, add:

- `evals/evals.json`
- small fixture files under `evals/files/`

Use evals to preserve the review behavior you care about.

## Step 7: Update repository documentation

After adding a skill, update:

- `README.md`
- `CHANGELOG.md`
- `docs/roadmap.md`

## Step 8: Test installation paths

At minimum, confirm one of these works:

- copy the skill directory back into `~/.claude/skills/`
- install from the Git repository with your preferred Agent Skills tool

## Step 9: Tag a release

Once the skill is stable, create a Git tag so others can install a known version.

## Checklist

- [ ] `SKILL.md` exists
- [ ] frontmatter is valid
- [ ] directory name matches `name`
- [ ] references use relative paths
- [ ] repo docs mention the new skill
- [ ] evals exist for subtle behaviors
