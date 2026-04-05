## Findings

- **Critical — worker goroutine never stops (`ctx` ignored), causing goroutine/memory leak**
  - `Start` receives `context.Context` but the goroutine loop never listens to `ctx.Done()`.

- **Critical — closed `jobs` channel triggers hot loop**
  - `case job := <-jobs:` does not check the `ok` flag.

## Recommended Fix Pattern

```go
func (p *Processor) Start(ctx context.Context, jobs <-chan string) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case job, ok := <-jobs:
				if !ok {
					return
				}
				p.handle(job)
			case <-ticker.C:
				p.flush()
			}
		}
	}()
}
```
