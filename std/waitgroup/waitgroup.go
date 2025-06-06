package waitgroup

import (
	"fmt"
	"github.com/ddkwork/golibrary/std/mylog"
	"sync"
)

type token struct{}

type Group struct {
	wg  sync.WaitGroup
	sem chan token
	mu  sync.RWMutex
}

func New() *Group {
	return &Group{
		wg:  sync.WaitGroup{},
		sem: nil,
		mu:  sync.RWMutex{},
	}
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}
	g.wg.Done()
}

func (g *Group) Wait() { //todo 其它地方使用select会导致死锁
	g.wg.Wait()
}

func (g *Group) Go(f func()) { //注意不要在回调里面使用log，应该在外部使用，否则会乱序，虽然内部加锁可以解决，但是影响效率
	if g.sem != nil {
		g.sem <- token{}
	}
	g.add(f)
}

func (g *Group) add(f func()) {
	g.wg.Add(1)
	go func() {
		defer g.done()
		mylog.Call(func() {
			//if g.UseMutex {
			//	g.mu.Lock() //action go get 无限等待，对于更新工作区下的模块同时多次读写，需要进一步测试，理论上需要加读写锁，那么应用场景只有这一个需要加锁？
			//	defer g.mu.Unlock()
			//}
			f()
		})
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
