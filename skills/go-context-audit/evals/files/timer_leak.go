package auditfixture

import (
	"context"
	"fmt"
	"time"
)

type EventBus struct {
	events chan string
}

func NewEventBus() *EventBus {
	return &EventBus{events: make(chan string, 100)}
}

func (b *EventBus) ConsumeLoop(ctx context.Context) {
	for {
		select {
		case evt := <-b.events:
			fmt.Println("event:", evt)
		case <-time.After(10 * time.Second):
			fmt.Println("idle timeout, flushing buffers")
		}
	}
}
