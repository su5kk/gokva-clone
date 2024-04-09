package stix

type SchedulerOption func(*Scheduler)

func WithRetryPipe(retrypipe chan string) SchedulerOption {
	return func(s *Scheduler) {
		s.retrypipe = retrypipe
	}
}
