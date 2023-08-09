package main

import (
	"context"
	"fmt"
	"time"

	"github.com/su5kk/stix"
)

func main() {
	// prepare data pipelines
	addpipe := make(chan *stix.Task)
	ctx := context.Background()

	scheduler := stix.NewScheduler(addpipe)
	shutdown := scheduler.Start(ctx)

	// submit task
	addpipe <- &stix.Task{
		ID:          "1",
		Name:        "my-task",
		Description: "print task",
		Execute: func() error {
			fmt.Println("HEllo world")
			return nil
		},
	}

	// wait for execution
	time.Sleep(1 * time.Second)

	// shutdown scheduler
	shutdown()
}
