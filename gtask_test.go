package golibs

import (
	"fmt"
	"testing"
	"time"
)

func tQueuedfun(event chan interface{}, a string, t time.Duration) {
	fmt.Println(a + " is Runnint")
	for {
		select {
		case evt := <-event:
			fmt.Println(a+",rcv:", evt)
		case <-time.After(t):
			fmt.Println(a + " Exit")
			return
		}
	}
}

func TestQueuedTask(t *testing.T) {
	jt := new(Tasks).Init(tQueuedfun, 1, 5, UNQUEUED_TASK)
	r := jt.TryRun(100*time.Millisecond, "FunA", 1*time.Second)
	fmt.Println("FunA ", r)
	r = jt.TryRun(100*time.Millisecond, "FunB", 3*time.Second)
	fmt.Println("FunB ", r)
	r = jt.TryRun(100*time.Millisecond, "FunC", 1*time.Second)
	fmt.Println("FunC ", r)
	r = jt.TryRun(100*time.Millisecond, "FunD", 3*time.Second)
	fmt.Println("FunD ", r)
	r = jt.TryRun(100*time.Millisecond, "FunE", 1*time.Second)
	fmt.Println("FunE ", r)

	time.Sleep(2 * time.Second)
	jt.Emit("test")

	time.Sleep(7 * time.Second)
}
