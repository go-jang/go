package util

import (
	"sort"
	"sync"
	"time"

	"github.com/go-jang/go/lang"
	"github.com/go-jang/go/util/concurrent"
)

type TimerTask struct {
	mutex             sync.Mutex
	cancelled         bool
	period            time.Duration
	nextExecutionTime time.Time
	run               func()
}

func NewTimerTask(run func()) *TimerTask {
	return &TimerTask{
		run: run,
	}
}

func (this *TimerTask) Cancel() bool {
	var result bool
	concurrent.Synchronized(&this.mutex, func() {
		result = !this.cancelled
		this.cancelled = true
	})
	return result
}

func (this *TimerTask) ScheduledExecutionTime() time.Time {
	var result time.Time
	concurrent.Synchronized(&this.mutex, func() {
		result = lang.If(this.period < 0, this.nextExecutionTime.Add(this.period), this.nextExecutionTime.Add(-this.period))
	})
	return result
}

type Timer struct {
	mutex     sync.Mutex
	cond      *sync.Cond
	cancelled bool
	queue     *ArrayList[*TimerTask]
}

func NewTimer() *Timer {
	t := Timer{}
	t.queue = NewArrayList[*TimerTask]()
	t.cond = sync.NewCond(&t.mutex)
	go t.mainLoop()
	return &t
}

func (this *Timer) ScheduleWithDelay(task *TimerTask, delay time.Duration) {
	lang.Assert(delay >= 0, "Negative delay")
	this.sched(task, time.Now().Add(delay), 0)
}

func (this *Timer) ScheduleWithDelayPeriod(task *TimerTask, delay, period time.Duration) {
	lang.Assert(delay >= 0, "Negative delay")
	lang.Assert(period > 0, "Non-positive period")
	this.sched(task, time.Now().Add(delay), -period)
}

func (this *Timer) ScheduleWithDelayPeriodFixedRate(task *TimerTask, delay, period time.Duration) {
	lang.Assert(delay >= 0, "Negative delay")
	lang.Assert(period > 0, "Non-positive period")
	this.sched(task, time.Now().Add(delay), period)
}

func (this *Timer) ScheduleAtTime(task *TimerTask, time time.Time) {
	this.sched(task, time, 0)
}

func (this *Timer) ScheduleAtTimeWithPeriod(task *TimerTask, time time.Time, period time.Duration) {
	lang.Assert(period > 0, "Non-positive period")
	this.sched(task, time, -period)
}

func (this *Timer) ScheduleAtTimeWithPeriodFixedRate(task *TimerTask, time time.Time, period time.Duration) {
	lang.Assert(period > 0, "Non-positive period")
	this.sched(task, time, period)
}

func (this *Timer) sched(task *TimerTask, time time.Time, period time.Duration) {
	concurrent.Synchronized(&this.mutex, func() {
		lang.Assert(!this.cancelled, "Timer is cancelled")
		lang.Assert(!task.cancelled, "Task already cancelled")
		task.period = period
		task.nextExecutionTime = time
		pos := sort.Search(this.queue.Size(), func(i int) bool {
			return this.queue.At(i).nextExecutionTime.After(task.nextExecutionTime)
		})
		this.queue.AddAt(pos, task)
		if pos == 0 {
			this.cond.Broadcast()
		}
	})
}

func (this *Timer) Cancel() {
	concurrent.Synchronized(&this.mutex, func() {
		this.cancelled = true
		this.queue.Clear()
		this.cond.Broadcast()
	})
}

func (this *Timer) mainLoop() {
	for !this.cancelled {
		var task *TimerTask
		var taskFired bool
		concurrent.SynchronizedWithCondition(&this.mutex, this.cond, func(cond *sync.Cond) {
			for this.queue.Empty() {
				cond.Wait()
				if this.cancelled {
					return
				}
			}
			var currentTime time.Time
			var executionTime time.Time
			task = this.queue.First()
			concurrent.Synchronized(&task.mutex, func() {
				if task.cancelled {
					this.queue.RemoveFirst()
					return
				}
				currentTime = time.Now()
				executionTime = task.nextExecutionTime
				taskFired = !executionTime.After(currentTime)
				if taskFired {
					this.queue.RemoveFirst()
					if task.period != 0 {
						task.nextExecutionTime = lang.If(task.period < 0, currentTime.Add(-task.period), executionTime.Add(task.period))
						pos := sort.Search(this.queue.Size(), func(i int) bool {
							return this.queue.At(i).nextExecutionTime.After(task.nextExecutionTime)
						})
						this.queue.AddAt(pos, task)
					}
				}
			})
			if !taskFired {
				concurrent.WaitWithTimeout(this.cond, executionTime.Sub(currentTime))
			}
		})
		if taskFired {
			task.run()
		}
	}
}
