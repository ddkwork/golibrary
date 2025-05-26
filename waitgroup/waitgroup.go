package waitgroup

import (
	"fmt"
	"sync"
)

type token struct{}

type Group struct {
	wg  sync.WaitGroup
	sem chan token
	mu  sync.Mutex
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}

func (g *Group) Wait() { g.wg.Wait() }

func (g *Group) Go(f func()) {
	if g.sem != nil {
		g.sem <- token{}
	}
	g.add(f)
}

func (g *Group) add(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.done()
		g.mu.Lock()
		defer g.mu.Unlock()
		f()
	}()
}

func (g *Group) TryGo(f func()) bool {
	if g.sem != nil {
		select {
		case g.sem <- token{}:
			// Note: this allows barging iff channels in general allow barging.
		default:
			return false
		}
	}
	g.add(f)
	return true
}

func (g *Group) SetLimit(n int) {
	if n < 0 {
		g.sem = nil
		return
	}
	if len(g.sem) != 0 {
		panic(fmt.Errorf("waitgroup: modify limit while %v goroutines in the group are still active", len(g.sem)))
	}
	g.sem = make(chan token, n)
}
