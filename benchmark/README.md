# Benchmarks

This folder stores **human- or tool-generated benchmark runs** for skills in this repository. Benchmarks help you quickly check:

- whether a skill is **triggered** when it should be (given the same user prompt and fixtures),
- whether the model output satisfies **concrete expectations** (not just “looks good”),
- how much the skill **improves** results versus the same setup **without** the skill.

Benchmarks are **not** a substitute for unit tests. They measure agent behavior under a specific model and prompt, so treat numbers as **signals**, not proofs.

## Layout

Benchmarks are grouped by skill and a numbered iteration (a frozen snapshot of methodology and prompts):

```text
benchmark/
└── <skill-name>/
    └── iteration-<N>/
        ├── benchmark.json          # machine-readable aggregate + per-run results
        ├── benchmark.md            # short human summary (optional but recommended)
        ├── review.html             # optional HTML report viewer
        └── <eval-slug>/            # one directory per eval case
            ├── eval_metadata.json  # prompt + assertion texts for this case
            ├── with_skill/
            │   ├── grading.json    # pass/fail per expectation + evidence
            │   └── outputs/        # raw artifacts (transcript, review, notes)
            └── without_skill/
                ├── grading.json
                └── outputs/
```

**Eval cases** should stay aligned with the skill’s canonical fixtures:

- `skills/<skill-name>/evals/evals.json` — eval definitions (ids, prompts, file paths, expectations)
- `skills/<skill-name>/evals/files/*.go` — code under review

The benchmark folder **mirrors** those evals as **slug** directories (for example `goroutine-cancel-leak` for the eval that uses `goroutine_leak.go`).

### Key files

| File | Role |
|------|------|
| `eval_metadata.json` | Frozen **prompt** and **assertions** for one eval; use this to rerun the same instruction text. |
| `with_skill` / `without_skill` | Same prompt and files; only difference should be whether the **skill** is enabled for the agent. |
| `grading.json` | Per-expectation `passed` / `failed` and short **evidence** quotes; defines “correct enough” for that run. |
| `benchmark.json` | All runs + `run_summary` (e.g. mean pass rate per configuration). |

## How to read a benchmark (no rerun)

1. Open `benchmark/<skill>/iteration-<N>/benchmark.md` for a one-page summary, or `benchmark.json` for detail.
2. For a single eval, open `with_skill/grading.json` and `without_skill/grading.json` and compare `summary.pass_rate` and each expectation.
3. Open `outputs/review.md` (and `transcript.md` if present) to see **why** a check passed or failed.

## How to reproduce “the same” benchmark

LLM outputs are **not bitwise-deterministic**. Reproducing means: **same inputs and procedure**, then **compare** pass rates and grading; expect small drift unless you fix the model version and use very low temperature.

Do the following:

1. **Pin the skill revision**  
   Use the same git commit (or tag) of this repo, or the same files under `skills/<skill-name>/`. Note the commit hash in your notes or in benchmark metadata if you publish a new iteration.

2. **Pin the agent stack**  
   Record at least: **client** (e.g. Claude Code), **model name**, and whether **tools** / **subagents** were used. The `benchmark.json` `metadata` block may include `executor_model` and `analyzer_model` for this purpose.

3. **Use the frozen prompt and assertions**  
   From `eval_metadata.json`, copy `prompt` into the user message. Attach or paste the same Go files as listed in `skills/<skill-name>/evals/evals.json` for that eval id (paths relative to the skill directory).

4. **Run two configurations**  
   - **with_skill**: skill installed and active (same trigger conditions you use in production).  
   - **without_skill**: same model and tools, but the skill disabled or not installed.

5. **Grade with the same rubric**  
   For each string in `assertions`, decide pass/fail against the model output and record evidence (as in `grading.json`). You can do this manually or with a separate grader model; keep the **assertion text** identical for comparability.

6. **Optional: reduce variance**  
   Set `runs_per_configuration` &gt; 1 in your procedure and aggregate mean/stddev (as in `benchmark.json` `run_summary`). Storing multiple runs per eval makes comparisons more stable.

## Adding a new benchmark iteration

When you change a skill or want to record a new model:

1. Copy the previous `iteration-<N>` tree to `iteration-<N+1>` (or create a fresh layout matching the convention above).
2. Update `benchmark.json` `metadata`: `timestamp`, models, `evals_run`, `runs_per_configuration`, and `skill_path` if relevant.
3. Rerun evals, refresh `grading.json` and `outputs/*`, then regenerate `benchmark.md` / `benchmark.json` summaries.
4. Mention the iteration in your changelog or PR so others know which snapshot is current.

## Relationship to `skills/.../evals/`

| Location | Purpose |
|----------|---------|
| `skills/<skill>/evals/` | **Source of truth** for eval definitions and `.go` fixtures. |
| `benchmark/<skill>/iteration-*` | **Recorded runs** and scores for a specific environment and time. |

Keeping eval definitions in the skill directory makes the skill self-contained; the benchmark directory proves how those evals behaved under a particular setup.

## Troubleshooting

- **Scores differ from `benchmark.json` on your machine** — Expected. Compare methodology (model, temperature, skill version, prompt wording).
- **Absolute paths in old artifacts** — Some saved paths may point to a local machine. For reproduction, use the repo’s `skills/.../evals/files/` paths and ignore host-specific prefixes when reading citations.
- **Zeros for time/tokens** — Some iterations omit executor metrics; see `notes` in `benchmark.json` or `benchmark.md`.
