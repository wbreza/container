package container

import (
	"context"
	"time"
)

// An emptyCtx is never canceled, has no values, and has no deadline.
// It is the common base of backgroundCtx and todoCtx.
type emptyCtx struct{}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyCtx) Done() <-chan struct{} {
	return nil
}

func (emptyCtx) Err() error {
	return nil
}

func (emptyCtx) Value(key any) any {
	return nil
}

func (emptyCtx) String() string {
	return "context.Empty"
}

func noContext() context.Context {
	return emptyCtx{}
}
