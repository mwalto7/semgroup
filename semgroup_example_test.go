package semgroup_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mwalto7/semgroup"
)

func ExampleGroup() {
	// Create a new group with a maximum of 10 in-flight goroutines.
	sg, ctx := semgroup.WithContext(context.Background(), 10)

	// At most 10 goroutines will print "subtask <i>" between sleeps.
	for i := 0; i < 100; i++ {
		i := i
		sg.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				fmt.Printf("subtask %d\n", i)
			}
			time.Sleep(2 * time.Second)
			return nil
		})
	}

	// Wait for all goroutines to finish, propagating the first non-nil error (if-any).
	if err := sg.Wait(); err != nil {
		log.Fatal(err)
	}
}
