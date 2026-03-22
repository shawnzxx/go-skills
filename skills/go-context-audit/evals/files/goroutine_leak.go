package auditfixture

import (
	"context"
	"time"
)

type Processor struct{}

func (p *Processor) Start(ctx context.Context, jobs <-chan string) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case job := <-jobs:
				p.handle(job)
			case <-ticker.C:
				p.flush()
			}
		}
	}()
}

func (p *Processor) handle(job string) {}

func (p *Processor) flush() {}
