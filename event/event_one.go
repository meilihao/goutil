package event

import (
	"context"
	"sync/atomic"
	"time"
)

// 事件驱动器
// 同一时刻仅运行一个任务
type EventOne struct {
	inputCh       chan int64
	num           int64                            // 正在进行的任务数
	remainNum     int64                            // 当前任务进行中, 又有任务进来
	workFn        func(context.Context, *EventOne) // need call EventOne.Done()
	isIgnoreOther bool                             // 是否忽略正在运行时, 进来的另一个任务
}

func NewEventOne(inputCap int, workFn func(context.Context, *EventOne), isIgnoreOther bool) *EventOne {
	if inputCap < 0 {
		inputCap = 0
	}

	return &EventOne{
		inputCh:       make(chan int64, inputCap),
		workFn:        workFn,
		isIgnoreOther: isIgnoreOther,
	}
}

func (eb *EventOne) Do(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-eb.inputCh:
				if atomic.LoadInt64(&eb.num) > 0 {
					atomic.AddInt64(&eb.remainNum, 1)

					continue
				}

				atomic.AddInt64(&eb.num, 1)
				go eb.workFn(ctx, eb)
			}
		}
	}()
}

func (eb *EventOne) Ticker(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()

				return
			case <-ticker.C:
				eb.Emit(0)
			}
		}
	}()
}

// 需注意一种情况: workFn中虽然使用了了Emit()但因为逻辑原因被跳过, 导致eb.inputCh一直没有传入而导致EventOne饿死
func (eb *EventOne) Emit(d time.Duration) {
	if int(d) == 0 {
		eb.inputCh <- time.Now().Unix()

		return
	}

	go func() {
		t := time.NewTimer(d)
		defer t.Stop()

		now := <-t.C
		eb.inputCh <- now.Unix()
	}()
}

func (eb *EventOne) Done() {
	atomic.AddInt64(&eb.num, -1)

	if eb.isIgnoreOther {
		atomic.StoreInt64(&eb.remainNum, 0)
	} else {
		if atomic.LoadInt64(&eb.remainNum) > 0 {
			atomic.StoreInt64(&eb.remainNum, 0)
			eb.Emit(0)
		}
	}
}

// func Work(ctx context.Context, eb *EventOne) {
// 	defer func() {
// 		eb.Done()
// 	}()

// 	fmt.Println(time.Now())
// 	time.Sleep(3 * time.Second)
// }

// func main() {
// 	entry := NewEventOne(10, Work, false)

// 	ctx, cancelFn := context.WithCancel(context.Background())
// 	defer cancelFn()

// 	entry.Do(ctx)

// 	fmt.Println("---", time.Now())
// 	entry.Emit(0)
// 	entry.Emit(2 * time.Second)
// 	entry.Emit(5 * time.Second)
// 	entry.Emit(4 * time.Second)

// 	select {}
// }
