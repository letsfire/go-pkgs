package limit

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSingle_Do(t *testing.T) {
	val, err := Single.Run("key", func() (interface{}, error) {
		return "value", nil
	})
	if val.(string) != "value" || err != nil {
		t.Errorf("unexpected result, expect = [value, nil], but = [%#v, %#v]", val, err)
	}
}

func TestSingle_DoErr(t *testing.T) {
	someErr := errors.New("some error")
	val, err := Single.Run("key", func() (interface{}, error) {
		return nil, someErr
	})
	if val != nil || err != someErr {
		t.Errorf("unexpected result, expect = [nil, %#v], but = [%#v, %#v]", someErr, val, err)
	}
}

func TestSingle_DoConcurrent(t *testing.T) {
	var id int32
	var wg sync.WaitGroup

	ch := make(chan string)
	cb := func() (interface{}, error) {
		atomic.AddInt32(&id, 1)
		return <-ch, nil
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Single.Run("key", cb)
			wg.Done()
		}()
	}

	// ensure that all goroutines are ready
	time.Sleep(50 * time.Millisecond)
	ch <- "value"
	wg.Wait()

	if atomic.LoadInt32(&id) != 1 {
		t.Error("unexpected result , should be equal to 1")
	}
}
