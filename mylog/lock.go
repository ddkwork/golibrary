package mylog

import (
	"fmt"
	"sync"
)

// go run -race  .
var (
	count int64
	lock  = new(sync.Mutex)
)

func main() {
	ch := make(chan struct{}, 2)
	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			count++
			lock.Unlock()
		}
		ch <- struct{}{}
	}()

	go func() {
		for i := 0; i < 100000; i++ {
			lock.Lock()
			count--
			lock.Unlock()
		}
		ch <- struct{}{}
	}()

	<-ch
	<-ch
	close(ch)
	fmt.Println(count)
}
