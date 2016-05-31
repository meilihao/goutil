# 基于Golang的异步任务包

from : https://github.com/lifei6671/async

## 使用

##### 示例1

```
	handler := func() string {
		fmt.Println("aaaaaaa");
		return "完成了";
	}
        //新建一个任务
	task := NewTask(handler);
        //运行并等待任务完成
	task.Run().Wait();
        //打印任务执行结果
	fmt.Println(task.Result)
```

##### 示例2

```
        handler1 := func(){
		fmt.Println(time.Now())
		fmt.Println("handler1","我在等待指定的时间后执行");
	}
	param2 := "aaaaaaaaaaaaaaaaaaaaaa";
	handler2 := func(p string) {
		fmt.Println("handler2",time.Now())
		fmt.Println(p);
	}
	param3 := "bbbbbbbbbbbbbbbbbbbbbbbbbb";
	handler3 := func(p string) string {
		fmt.Println(p);
		return p+"111111111111111";
	}

        //新建一个任务并添加执行完成后的延续任务，同时指定任务延时执行时间
	task1 := NewTask(handler1).ContinueWith(func(result TaskResult){
		fmt.Println("我在task1执行后执行。");
	}).ContinueWith(func(result TaskResult){
		fmt.Println("我在task1执行后第二次执行。");
	}).Delay(5*time.Second).Run();
	task2 := NewTask(handler2,param2).Run();
	task3 := NewTask(handler3,param3).Run();

        //批量等待任务完成
	WaitAll(task1,task2,task3);
	fmt.Println(task3.Result)
```

##### 示例3

```
	handler1 := func(){
		fmt.Println(time.Now())
		fmt.Println("TestTaskList_Add: handler1");
	}
	param2 := "TestTaskList_Add:aaaaaaaaaaaaaaaaaaaaaa";
	handler2 := func(p string) {
		fmt.Println("TestTaskList_Add:handler2",time.Now())
		fmt.Println(p);
	}
	param3 := "TestTaskList_Add: bbbbbbbbbbbbbbbbbbbbbbbbbb";
	handler3 := func(p string) string {
		fmt.Println(p);
		return p+"111111111111111";
	}

	task1 := NewTask(handler1);
	task2 := NewTask(handler2,param2);
	task3 := NewTask(handler3,param3);

        //新建一个任务列表
	taskList := NewTaskList();

        //添加任务到列表中，同时执行任务，然后等待所有任务完成。
	taskList.AddRange(task1,task2,task3).Run().WaitAll();
```
