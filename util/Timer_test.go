package util_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-jang/go/util"
)

func Test_Timer_ScheduleWithDelayPeriod(t *testing.T) {
	t.Run("subsribtion test'", func(t *testing.T) {
		t.Skip("for manual run")
		timer := util.NewTimer()
		timer.ScheduleWithDelayPeriod(util.NewTimerTask(func() {
			fmt.Println(time.Now())
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		}), time.Second, time.Second)
		time.Sleep(15 * time.Second)
	})
}

func Test_Timer_ScheduleWithDelayPeriodFixedRate(t *testing.T) {
	t.Run("subsribtion test'", func(t *testing.T) {
		t.Skip("for manual run")
		timer := util.NewTimer()
		timer.ScheduleWithDelayPeriodFixedRate(util.NewTimerTask(func() {
			fmt.Println(time.Now())
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		}), time.Second, time.Second)
		time.Sleep(15 * time.Second)
	})
}

func Test_Timer_ScheduleAtFixedRateAtTimeWithPeriod(t *testing.T) {
	t.Run("subsribtion test'", func(t *testing.T) {
		t.Skip("for manual run")
		moment := time.Now().Add(10 * time.Second)
		timer := util.NewTimer()
		timer.ScheduleAtTime(util.NewTimerTask(func() {
			fmt.Println(time.Now())
		}), moment)
		time.Sleep(15 * time.Second)
	})
}
