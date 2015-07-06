package gracefulshutdown

import (
	"errors"
	"testing"
	"time"
)

type SMPingFunc func()

func (f SMPingFunc) Ping() {
	f()
}

func (f SMPingFunc) ShutdownFinish() {

}

func (f SMPingFunc) Start(ssi StartShutdownInterface) error {
	return nil
}

type SMFinishFunc func()

func (f SMFinishFunc) Ping() {

}

func (f SMFinishFunc) ShutdownFinish() {
	f()
}

func (f SMFinishFunc) Start(ssi StartShutdownInterface) error {
	return nil
}

type SMStartFunc func() error

func (f SMStartFunc) Ping() {

}

func (f SMStartFunc) ShutdownFinish() {

}

func (f SMStartFunc) Start(ssi StartShutdownInterface) error {
	return f()
}

func TestCallbacksGetCalled(t *testing.T) {
	gs := New(time.Millisecond)

	c := make(chan int, 100)
	for i := 0; i < 15; i++ {
		gs.AddShutdownCallback(ShutdownFunc(func() error {
			c <- 1
			return nil
		}))
	}

	gs.StartShutdown(SMPingFunc(func() {}))

	if len(c) != 15 {
		t.Error("Expected 15 elements in channel, got ", len(c))
	}
}

func TestStartGetsCalled(t *testing.T) {
	gs := New(time.Hour)

	c := make(chan int, 100)
	for i := 0; i < 15; i++ {
		gs.AddShutdownManager(SMStartFunc(func() error {
			c <- 1
			return nil
		}))
	}

	gs.Start()

	if len(c) != 15 {
		t.Error("Expected 15 Start to be called, got ", len(c))
	}
}

func TestStartErrorGetsReturned(t *testing.T) {
	gs := New(time.Hour)

	gs.AddShutdownManager(SMStartFunc(func() error {
		return errors.New("my-error")
	}))

	err := gs.Start()
	if err == nil || err.Error() != "my-error" {
		t.Error("Shutdown did not return my-error, got ", err)
	}
}

func TestPingGetsCalled(t *testing.T) {
	c := make(chan int, 100)
	gs := New(2 * time.Millisecond)

	gs.AddShutdownCallback(ShutdownFunc(func() error {
		time.Sleep(5 * time.Millisecond)
		return nil
	}))

	gs.StartShutdown(SMPingFunc(func() {
		c <- 1
	}))

	time.Sleep(5 * time.Millisecond)

	if len(c) != 3 {
		t.Error("Expected 3 pings, got ", len(c))
	}
}

func TestShutdownFinishGetsCalled(t *testing.T) {
	c := make(chan int, 100)
	gs := New(2 * time.Millisecond)

	gs.AddShutdownCallback(ShutdownFunc(func() error {
		time.Sleep(5 * time.Millisecond)
		return nil
	}))

	gs.StartShutdown(SMFinishFunc(func() {
		c <- 1
	}))

	if len(c) != 1 {
		t.Error("Expected 1 ShutdownFinish, got ", len(c))
	}
}
