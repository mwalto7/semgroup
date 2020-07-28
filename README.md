# semgroup

[![GoDoc](https://godoc.org/github.com/mwalto7/semgroup?status.svg)](https://pkg.go.dev/github.com/mwalto7/semgroup?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/mwalto7/semgroup)](https://goreportcard.com/report/github.com/mwalto7/semgroup)

`semgroup` provides a simple wrapper around an [error group](https://godoc.org/golang.org/x/sync/errgroup)
that adds the ability to limit the maximum number of in-flight goroutines working on a group of tasks
via a [weighted semaphore](https://godoc.org/golang.org/x/sync/semaphore). The API is exactly the same
as the `errgroup` package.

## Getting Started

Get the latest version of the package with `go get`

```
go get -u github.com/mwalto7/semgroup
```

## Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "runtime"
    "time"

    "github.com/mwalto7/semgroup"
)

func main() {
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
```
