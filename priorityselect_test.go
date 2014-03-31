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
		t.Log(msg)
		if msg != expect {
			t.Fatalf("expect:", expect)
		}
		time.Sleep(time.Millisecond * 200)
	}
}
