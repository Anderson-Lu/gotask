[Chinese Document](https://github.com/Anderson-Lu/gotask/blob/master/readme_cn.md)

Gotask is a concurrent task scheduling tool based on waitegroup. Concurrent task solution for concurrent execution of computation and consistency of result data

### Dependency

```shell
go get github.com/Anderson-Lu/gotask
```

### Quick Start

Create a new task manager like this:

```golang
var tasks = NewGoTask(500, false)
```

here,`500` is the largest concurrent number for `taksManager`. if `quickMode` is set to true, specifc task will be executed immediately when `task.AddTask()` method is invoked.

Defined a function witch will be excuted concurrently. Note that all parameters is `...interface{}` and use `tasks.GetParamter` to get concrete parameter.

```golang
var total = 0
func add(params ...interface{}) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	lock.Lock()
	total += tasks.GetParamter(0, params).(int)
	lock.Unlock()
}
```

Regist concret task to manager.

```golang
for i := 0; i < 10000; i++ {
	tasks.Add(add, 1)
}
tasks.Start() 
```

if `tasks.Start()` is invoked, main routine will be blocked and wait to all subTask finish.

### Sample Demo

Cyclic 10000 cumulative calculation

```golang
var tasks = NewGoTask(10000, false)
var lock sync.Mutex
var total = 0

func add(params ...interface{}) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	lock.Lock()
	total += tasks.GetParamter(0, params).(int)
	lock.Unlock()
}

func Demo() {
	for i := 0; i < 10000; i++ {
		tasks.Add(add, 1)
	}
	t.StartTimer()
	tasks.Start()
}
```


### API Intro

```shell
//Create concurrent task manager instance
NewGoTask(maxConcurentNum int, quickMode bool) *GoTask 

//Add subtask
Add(task func(...interface{}), params ...interface{})

//Get specific params with index
GetParamter(index int, params interface{}) interface{} 

//Start all subTasks concurrently
Start()

//Return how much time cost for all subtasks
Cost() int

//Finsh all tasks manually and is required if quickmode = true
Done()
```

### Performance

we assume that we add numbers in concurrent way. Everytime we add 1 and loop 10000 times.

```shell
conccurent_num    loop_time     cost/op            bytes/op       allocs/op         total_time  avg_time
--------------------------------------------------------------------------------------------------------
100               1             7816076500 ns/op   2476760 B/op   40396 allocs/op   7.82s       7.82s
500               1             1576902314 ns/op   2669496 B/op   41186 allocs/op   1.528s      1.528s
1000              2             791744067 ns/op    1668632 B/op   30433 allocs/op   2.420s      1.21s
10000             300000        3564 ns/op         37 B/op        0 allocs/op       7.841s      0.0000261s
```