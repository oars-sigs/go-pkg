package flow

import "sync"

type gawait struct {
	m sync.Map
}

func newAwait() *gawait {
	return &gawait{
		m: sync.Map{},
	}
}

func (w *gawait) Await(ss []string, ignore bool) error {
	for _, s := range ss {
		v, ok := w.m.Load(s)
		if !ok {
			continue
		}
		err := <-v.(chan error)

		if err != nil && !ignore {
			return err
		}

	}
	return nil
}

func (w *gawait) AddAwait(s string) {
	a := make(chan error)
	w.m.Store(s, a)
}

func (w *gawait) DoneAwait(s string, err error) {
	v, ok := w.m.Load(s)
	if !ok {
		return
	}
	go func() {
		v.(chan error) <- err
	}()
}
