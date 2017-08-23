package cflib

import (
	"time"
)

type Timer interface {
	SleepSec(seconds int)
	Stop()
}

// Time handle Gos time object
type Time struct {
	timer  *time.Timer
	cancel chan interface{}
}

// Sleep in seconds
func (t *Time) SleepSec(seconds int) {
	t.timer = time.NewTimer(time.Duration(seconds) * time.Second)
	t.cancel = make(chan interface{})
	for {
		select {
		case <-t.timer.C:
			return
		case <-t.cancel:
			return
		}
	}
}

// Stop running timer
func (t *Time) Stop() {
	if t.timer != nil {
		t.timer.Stop()
		close(t.cancel)
	}
}

// TimeStub stubs the time object
type TimeStub struct{}

func (t *TimeStub) SleepSec(seconds int) {

}

func (t *TimeStub) Stop() {

}
