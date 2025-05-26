package waitgroup_test

import (
	"fmt"
	"mcp/waitgroup"
	"sync/atomic"
	"testing"
	"time"
)

func TestTryGo(t *testing.T) {
	g := &waitgroup.Group{}
	n := 42
	g.SetLimit(42)
	ch := make(chan struct{})
	fn := func() {
		ch <- struct{}{}
	}
	for i := range n {
		if !g.TryGo(fn) {
			t.Fatalf("TryGo should succeed but got fail at %d-th call.", i)
		}
	}
	if g.TryGo(fn) {
		t.Fatalf("TryGo is expected to fail but succeeded.")
	}
	go func() {
		for range n {
			<-ch
		}
	}()
	g.Wait()

	if !g.TryGo(fn) {
		t.Fatalf("TryGo should success but got fail after all goroutines.")
	}
	go func() { <-ch }()
	g.Wait()

	// Switch limit.
	g.SetLimit(1)
	if !g.TryGo(fn) {
		t.Fatalf("TryGo should success but got failed.")
	}
	if g.TryGo(fn) {
		t.Fatalf("TryGo should fail but succeeded.")
	}
	go func() { <-ch }()
	g.Wait()

	// Block all calls.
	g.SetLimit(0)
	for range 1 << 10 {
		if g.TryGo(fn) {
			t.Fatalf("TryGo should fail but got succeded.")
		}
	}
	g.Wait()
}

func TestGoLimit(t *testing.T) {
	const limit = 10

	g := &waitgroup.Group{}
	g.SetLimit(limit)
	var active int32
	for i := 0; i <= 1<<10; i++ {
		g.Go(func() {
			n := atomic.AddInt32(&active, 1)
			if n > limit {
				panic(fmt.Errorf("saw %d active goroutines; want ≤ %d", n, limit))
			}
			time.Sleep(1 * time.Microsecond) // Give other goroutines a chance to increment active.
			atomic.AddInt32(&active, -1)
		})
	}
	g.Wait()
}

func BenchmarkGo(b *testing.B) {
	fn := func() {}
	g := &waitgroup.Group{}
	b.ReportAllocs()
	for b.Loop() {
		g.Go(func() { fn() })
	}
	g.Wait()
}
