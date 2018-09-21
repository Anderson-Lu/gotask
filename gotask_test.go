package gotask

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

var tasks = NewGoTask(500, false)
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
