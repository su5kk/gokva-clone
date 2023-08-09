package main

import (
	"context"
	"fmt"

	"github.com/su5kk/stix"
)

func main() {
	addpipe := make(chan *stix.Task)
	retrypipe := make(chan string)
	scheduler := stix.NewScheduler(addpipe, retrypipe)
	ctx := context.Background()
	shutdown := scheduler.Start(ctx)
	addpipe <- &stix.Task{
		ID:          "1",
		Name:        "my-task",
		Description: "print task",
		Execute: func() error {
			fmt.Println("HEllo world")
			return nil
		},
	}
	shutdown()
}
