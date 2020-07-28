package semgroup_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/mwalto7/semgroup"
)

func TestGroup_Go(t *testing.T) {
	start := func(ctx context.Context, id, max int) func() error {
		// Save number of goroutines running before starting subtasks.
		preExisting := runtime.NumGoroutine()
		return func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			// Ensure there are not >g.Weight() in-flight goroutines at any given time.
			if n := runtime.NumGoroutine() - preExisting; n > max {
				return fmt.Errorf("subtask %d: want at most %d goroutines, got %d", id, max, n)
			}
			return nil
		}
	}

	g, ctx := semgroup.WithContext(context.Background(), -1)

	// Start g.Weight()*2 subtasks.
	weight := int(g.Weight())
	for id := 0; id < weight*2; id++ {
		g.Go(start(ctx, id, weight))
	}

	// Fail the test if there were ever more than g.Weight() in-flight goroutines.
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}
