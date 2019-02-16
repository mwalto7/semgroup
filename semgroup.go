// Package semgroup provides a modified errgroup.Group with the ability to limit
// the maximum concurrent access of resources for groups of goroutines working on
// subtasks of a common task.
package semgroup

import (
	"context"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// Group is an errgroup.Group combined with a semaphore. The
// maximum number of in-flight goroutines is equal to
// the weight of the semaphore.
//
// A zero Group is invalid. Use WithContext to initialize a Group.
type Group struct {
	ctx context.Context
	eg  *errgroup.Group
	sem *semaphore.Weighted
}

// WithContext returns a new Group with a weighted semaphore with
// weight n and an associated Context derived from ctx.
//
// If the given Context is nil, context.Background is used.
//
// If the given semaphore weight n is less than or equal to zero,
// runtime.NumCPU()*2 is used.
func WithContext(ctx context.Context, n int64) (*Group, context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	if n <= 0 {
		n = int64(runtime.NumCPU() * 2)
	}
	sg := &Group{sem: semaphore.NewWeighted(n)}
	sg.eg, sg.ctx = errgroup.WithContext(ctx)
	return sg, sg.ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (sg *Group) Wait() error {
	return sg.eg.Wait()
}

// Go blocks until a semaphore is acquired, then calls the given function in
// a new goroutine. If a non-nil error is returned when acquiring the
// semaphore, the error is returned and the group is cancelled.
//
// See errgroup for details on error handling and propagation.
func (sg *Group) Go(f func() error) {
	err := sg.sem.Acquire(sg.ctx, 1)
	sg.eg.Go(func() error {
		if err != nil {
			return err
		}
		defer sg.sem.Release(1)

		return f()
	})
}
