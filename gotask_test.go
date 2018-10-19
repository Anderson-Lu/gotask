package gotask

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var tasks = NewGoTask(10000, false)
var lock sync.Mutex
var total = 0

func add(params ...interface{}) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	lock.Lock()
	total += tasks.GetParamter(0, params).(int)
	lock.Unlock()
}

func BenchmarkConcureentTask(t *testing.B) {
	for i := 0; i < 10000; i++ {
		tasks.Add(add, 1)
	}
	t.StartTimer()
	tasks.Start() 
}

func BenchmarkTypeAssert(t *testing.B) {
	t.StartTimer()
	var wg sync.WaitGroup
	var k interface{} = "s"
	sum := 0
	for i := 0; i < 10000; i++ {
		go func(s interface{}) {
			defer func() {
				wg.Done()
			}()
			wg.Add(1)
			sum += len(k.(string))
		}(k)
	}
	wg.Wait()
}

//basic usage via gorutine but specific concreate param type
func BenchmarkConcurrentTaskRoutine(t *testing.B) {
	t.StartTimer()
	sum := 0
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		go func(i int) {
			defer func() {
				wg.Done()
			}()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			sum += 1
			wg.Add(1)
		}(i)
	}
	wg.Wait()
}
