package stix_test

import (
	"context"
	"testing"

	"github.com/su5kk/stix"
)

func TestScheduler(t *testing.T) {
	resultpipe := make(chan string)
	addpipe := make(chan *stix.Task)
	scheduler := stix.NewScheduler(addpipe)
	shutdown := scheduler.Start(context.Background())
	defer shutdown()

	addpipe <- &stix.Task{
		ID:          "id-1",
		Name:        "name-1",
		Description: "desc-1",
		Execute: func(context.Context) error {
			resultpipe <- "result"
			return nil
		},
	}

	result := <-resultpipe
	if result != "result" {
		t.Fail()
	}
}
