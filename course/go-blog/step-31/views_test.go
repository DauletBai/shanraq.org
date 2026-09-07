package main

import (
	"sync"
	"testing"
)

// TestViewsUnderLoad is the test that has to be run with -race. Without the
// mutex in views the detector names the line; with it the report is empty and
// the count comes out exact.
func TestViewsUnderLoad(t *testing.T) {
	v := newViews()

	const readers = 200

	var wg sync.WaitGroup
	for i := 0; i < readers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v.Add("dala")
		}()
	}
	wg.Wait()

	if got := v.Count("dala"); got != readers {
		t.Errorf("оқылған саны %d, күткеніміз %d", got, readers)
	}
}

func TestViewsTotalIsACopy(t *testing.T) {
	v := newViews()
	v.Add("dala")

	total := v.Total()
	total["dala"] = 100

	if got := v.Count("dala"); got != 1 {
		t.Errorf("көшірмені өзгерту санауышқа тиді: %d", got)
	}
}
