package priorityselect

import (
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)
	selector := New(c1, c2)

	go func() {
		time.Sleep(time.Millisecond * 50)
		for i := 0; i < 3; i++ {
			c1 <- "one"
			time.Sleep(time.Millisecond * 50)
		}
	}()

	go func() {
		for i := 0; i < 3; i++ {
			c2 <- "two"
			time.Sleep(time.Millisecond * 50)
		}
	}()

	expects := []string{"two", "one", "one", "one", "two", "two"}

	for _, expect := range expects {
		msg, err := selector.Select()
		if err != nil {
			t.Fatal(err)
		}
		if msg != expect {
			t.Fatalf("expect:", expect)
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func TestSelect2(t *testing.T) {
	c := make(chan int)
	n := 1024 * 1024
	go func() {
		for i := 0; i < n; i++ {
			c <- i
		}
	}()
	selector := New(c)
	for i := 0; i < n; i++ {
		_, err := selector.Select()
		if err != nil {
			t.Fail()
		}
	}
}

func TestSelect3(t *testing.T) {
	nChans := 128
	nMsgs := 128
	chans := make([]interface{}, 0)
	for i := 0; i < nChans; i++ {
		c := make(chan int)
		chans = append(chans, c)
		go func(n int) {
			for i := 0; i < nMsgs; i++ {
				c <- n*nChans + i
			}
		}(i)
	}
	selector := New(chans...)
	for i := 0; i < nChans*nMsgs; i++ {
		_, err := selector.Select()
		if err != nil {
			t.Fail()
		}
	}
}
