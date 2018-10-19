[English Document](https://github.com/Anderson-Lu/gotask/blob/master/readme.md)

Gotask是一个基于waitgroup实现的并发任务调度工具。用于并发执行计算并将保证结果数据一致性的并发任务解决方案。

### 依赖

```shell
go get github.com/Anderson-Lu/gotask
```

### 快速开始

创建一个并发任务管理实例:

```golang
var tasks = NewGoTask(500, false)
```

这里，`500`是最大的并发数。如果将第二个参数设置为`true`，则在调用`Adtask`方法时，将立即执行指定任务。

接下来定义一个同时被运行的函数。注意，所有参数都是`…interface{}`，并使用`GetParameter()`来获取具体参数。

```golang
var total = 0
func add(params ...interface{}) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	lock.Lock()
	total += tasks.GetParamter(0, params).(int)
	lock.Unlock()
}
```

向并发任务管理器注册任务:

```golang
for i := 0; i < 10000; i++ {
	tasks.Add(add, 1)
}
tasks.Start() 
```

如果调用`Stant()`，主例程将被阻塞并等待所有子任务完成。


### 简单示例

循环10000次累加计算

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


### API介绍

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

### 性能

以下是并发累加10000次的运行比较:

```shell
conccurent_num    loop_time     cost/op            bytes/op       allocs/op         total_time  avg_time
--------------------------------------------------------------------------------------------------------
100               1             7816076500 ns/op   2476760 B/op   40396 allocs/op   7.82s       7.82s
500               1             1576902314 ns/op   2669496 B/op   41186 allocs/op   1.528s      1.528s
1000              2             791744067 ns/op    1668632 B/op   30433 allocs/op   2.420s      1.21s
10000             300000        3564 ns/op         37 B/op        0 allocs/op       7.841s      0.0000261s
```

#### 有一些原因使得GoStand慢于原来的goroutine：

- 在我的最后一次提交中，我在每个任务之间添加了一个时间分割持续时间（50ms）,最新提交已删除。

```golang
for {
	if self.curTaskNum < self.max {
		break
	}
	// time.Sleep(time.Millisecond * 50) //now I delete this row 
}
```

- 如果你想使并发任务更快，你可以像这样设置`quickMode=true`：

```golang
var tasks = NewGoTask(10000, true)
```

- 对于某些情况，如果您的参数是具体的，我建议您使用特定的`gorutine`而不是使用`gotask`，因为这可以避免多类型转换。但是如果任务不同（如三个任务，taskA、taskB和taskC），那么可以使用`gotask`，因为它使代码更加优雅和易于维护。`gotask_test.go`中的demo是简单示例，所有任务都是相同的，因此如果使用原始Grututin而不是Gotask。


那么，是什么造就了Gotask慢起来呢？

- `params`类型转换(`interface{}`到具体的数据类型)
- 外部管理成本(比如耗时统计)
-所有任务都是闭包方法

那么，何时使用Gotask？

- 多个任务是不同的,即多个任务处理程序.
- 协程数量限制，尤其是在一些爬虫环境下，太大的并发数量将使你的机器瘫痪。
- 你需要统计所有任务的执行情况。