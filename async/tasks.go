package tasks

import (
	"sync"
	"reflect"
	"log"
	"container/list"
	"time"
)


//参数
type TaskParameter interface{};

//执行的方法
type TaskHanlder interface {};

//等待任务执行完成的后续任务
type ContinueWithHandler func(TaskResult);

//返回的参数类型
type TaskResult struct {
	Result interface{}
	Error error
}
//一个任务
type Task struct {
	wait *sync.WaitGroup
	handler reflect.Value
	params []reflect.Value
	Result TaskResult	//任务执行完成的返回结果
	once sync.Once
	IsCompleted bool	//表示任务是否执行完成
	continueWith *list.List
	delay time.Duration
}
//新建一个任务
func NewTask(handler TaskHanlder,params ...TaskParameter) *Task {

	handlerValue := reflect.ValueOf(handler);

	if(handlerValue.Kind() == reflect.Func){
		task := Task{
			wait : &sync.WaitGroup{},
			handler : handlerValue ,
			IsCompleted:false,
			continueWith : list.New(),
			delay: 0*time.Second,
			params : make([]reflect.Value,0),
		}
		if paramNum := len(params);paramNum > 0{
			task.params = make([]reflect.Value,paramNum);
			for index, v := range params {
				log.Println(index);
				task.params[index] = reflect.ValueOf(v);
			}
		}
		return &task;
	}
	panic("handler not func");
}

//启动异步任务
func (task *Task) Run() *Task {
	task.once.Do(func() {
		task.wait.Add(1);
		if(task.delay.Nanoseconds() > 0){
			time.Sleep(task.delay);
		}

		go func(){
			defer func(){

				task.IsCompleted = true;
				if(task.continueWith != nil){
					result := task.Result;
					for element:= task.continueWith.Back();element != nil;element = element.Prev(){
						if tt,ok := element.Value.(ContinueWithHandler);ok{
							tt(result);
						}

					}
				}
				task.wait.Done();
			}()


			values := task.handler.Call(task.params);

			task.Result = TaskResult{
				Result : values,
			};
		}();
	});
	return task;
}

//等待任务完成
func(task *Task) Wait() {
	task.wait.Wait();
}
//等待所有任务都完成
func WaitAll(tasks ...*Task){
	wait := &sync.WaitGroup{};
	for _,task := range tasks{
		wait.Add(1);
		go func() {
			defer wait.Done();
			task.wait.Wait();
		}();
	}
	wait.Wait();
}
//立即启动一个异步任务
func StartNew(handler TaskHanlder,params ...TaskParameter) *Task {
	task := NewTask(handler,params);
	task.Run();
	return task;
}
//当前Task执行完后执行
func (task *Task) ContinueWith(handler ContinueWithHandler) *Task {

	task.continueWith.PushFront(handler);

	return task;
}
//延迟指定的时间后执行
func (task *Task)Delay (delay time.Duration) *Task {
	task.delay = delay;
	return task;
}