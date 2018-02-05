package tachikoma

import "fmt"

type Scraper struct {
	emitters []Emitter
	jobs     []Job
	input    chan interface{}
}

func (s *Scraper) PushEmitterFn(fn EmitterFnProto) {
	s.emitters = append(s.emitters, NewEmitterFn(fn))
}

func (s *Scraper) PushJobs(jobs ...Job) {
	s.jobs = append(s.jobs, jobs...)
}

func (s *Scraper) JobByName(name string) (Job, error) {
	for _, j := range s.jobs {
		if j.GetName() == name {
			return j, nil
		}
	}
	return nil, fmt.Errorf("Unknown job %s", name)
}

func (s Scraper) Start() {
	go s.start()
	s.input = make(chan interface{}, 100)
	for _, e := range s.emitters {
		go e.Start(s.input)
	}
}

func (s Scraper) start() {
	for i := range s.input {
		for _, j := range s.jobs {
			j.Run(i)
		}
	}
}
