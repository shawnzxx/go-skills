# Benchmarks

## Rule: outputs must come from real LLM runs

Every `iteration-<N>` must be produced by **actual model calls**. Never copy from previous iterations or write template content.

For every eval case you **must**:

1. Send the real `prompt` + fixture file to the LLM with the skill **enabled** → save output to `with_skill/outputs/review.md`
2. Send the same prompt + file with the skill **disabled** → save output to `without_skill/outputs/review.md`
3. Grade each output against the `assertions` in `eval_metadata.json` → save to `with_skill/grading.json` and `without_skill/grading.json`
4. Record model name, date, and whether skill was enabled in `outputs/transcript.md`

## Layout

```text
benchmark/
└── <skill-name>/
    └── iteration-<N>/
        ├── benchmark.json          # aggregate results (all runs + run_summary)
        ├── benchmark.md            # one-paragraph human summary
        └── <eval-slug>/
            ├── eval_metadata.json  # frozen prompt + assertion texts
            ├── with_skill/
            │   ├── grading.json    # pass/fail per assertion + evidence
            │   └── outputs/
            │       ├── review.md   # raw model output
            │       └── transcript.md  # model, date, skill on/off
            └── without_skill/
                ├── grading.json
                └── outputs/
```

## Source of truth

Eval definitions and fixture files live in the skill directory:

- `skills/<skill-name>/evals/evals.json` — prompts, file paths, assertions
- `skills/<skill-name>/evals/files/` — Go fixtures

The benchmark directory records what happened when those evals were run.

## Adding a new iteration

1. Create `iteration-<N+1>/` following the layout above.
2. Read `skills/<skill-name>/evals/evals.json` to get the prompts, assertions, and fixture paths.
3. For each eval: run `with_skill` and `without_skill` calls and save `review.md`.
4. Grade outputs against `assertions` and write `grading.json`.
5. Write `benchmark.json` (`runs` + `run_summary` with mean pass rates) and `benchmark.md` (summary + observations).

## grading.json schema

```json
{
  "expectations": [
    {
      "text": "<assertion>",
      "passed": true,
      "evidence": "<short quote from review>"
    }
  ],
  "summary": { "passed": 3, "failed": 1, "total": 4, "pass_rate": 0.75 },
  "execution_metrics": {},
  "timing": {}
}
```

## eval_metadata.json schema

```json
{
  "eval_id": 1,
  "eval_name": "<slug>",
  "prompt": "<user message sent to model>",
  "assertions": ["<expectation 1>", "<expectation 2>"]
}
```
