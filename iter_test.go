package sortedmap

import (
	"testing"
	"time"
)

const (
	nilBoundValsErr = "accepted nil bound values"
	generalBoundsErr = "between bounds error"
	nilValErr = "nil value!"
	nonNilValErr = "non-nil value"
	runawayIterErr = "runaway iterator!"
)

var maxTime = time.Unix(1<<63 - 62135596801, 999999999)

func TestIterCh(t *testing.T) {
	if _, _, err := newSortedMapFromRandRecords(1000); err != nil {
		t.Fatal(err)
	}
}

func TestIterChTimeout(t *testing.T) {
	sm, _, err := newSortedMapFromRandRecords(1000)
	if err != nil {
		t.Fatal(err)
	}
	timeout := 1 * time.Microsecond
	sleepDur := 10 * time.Millisecond
	params := IterChParams{SendTimeout: timeout}

	if ch, ok := sm.CustomIterCh(params); ok {
		for i := 0; i < 5; i++ {
			time.Sleep(sleepDur)
			rec := <- ch
			if i > 1 && rec.Key != nil {
				t.Fatalf("TestIterChTimeout failed: %v: %v", nonNilValErr, rec.Key)
			}
		}
	} else {
		t.Fatal("TestIterChTimeout failed: %v", generalBoundsErr)
	}

	params.LowerBound = time.Time{}
	params.UpperBound = maxTime

	if ch, ok := sm.CustomIterCh(params); ok {
		for i := 0; i < 5; i++ {
			time.Sleep(sleepDur)
			rec := <- ch
			if i > 1 && rec.Key != nil {
				t.Fatalf("TestIterChTimeout failed: %v: %v", nonNilValErr, rec.Key)
			}
		}
	} else {
		t.Fatal("TestIterChTimeout failed: %v", generalBoundsErr)
	}

	params = IterChParams{
		Reversed: true,
		SendTimeout: timeout,
	}
	if ch, ok := sm.CustomIterCh(params); ok {
		for i := 0; i < 5; i++ {
			time.Sleep(sleepDur)
			rec := <- ch
			if i > 1 && rec.Key != nil {
				t.Fatalf("TestIterChTimeout failed: %v: %v", nonNilValErr, rec.Key)
			}
		}
	} else {
		t.Fatal("TestIterChTimeout failed: %v", generalBoundsErr)
	}

	params.LowerBound = time.Time{}
	params.UpperBound = maxTime

	if ch, ok := sm.CustomIterCh(params); ok {
		for i := 0; i < 5; i++ {
			time.Sleep(sleepDur)
			rec := <- ch
			if i > 1 && rec.Key != nil {
				t.Fatalf("TestIterChTimeout failed: %v: %v", nonNilValErr, rec.Key)
			}
		}
	} else {
		t.Fatal("TestIterChTimeout failed: %v", generalBoundsErr)
	}
}

func TestBoundedIterCh(t *testing.T) {
	sm, _, err := newSortedMapFromRandRecords(1000)
	if err != nil {
		t.Fatal(err)
	}

	reversed := false
	earlierDate := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	laterDate := time.Date(5321, 1, 1, 0, 0, 0, 0, time.UTC)

	if ch, ok := sm.BoundedIterCh(reversed, nil, nil); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if ch, ok := sm.BoundedIterCh(reversed, time.Time{}, maxTime); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if ch, ok := sm.BoundedIterCh(reversed, maxTime, time.Time{}); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if ch, ok := sm.BoundedIterCh(reversed, earlierDate, time.Now()); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if ch, ok := sm.BoundedIterCh(reversed, time.Now(), earlierDate); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if ch, ok := sm.BoundedIterCh(reversed, time.Now(), laterDate); ok {
		if err := verifyRecords(ch, reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}

	if _, ok := sm.BoundedIterCh(reversed, laterDate, laterDate); ok {
		t.Fatal("TestBoundedIterCh failed: %v", generalBoundsErr)
	}
}

func TestCustomIterCh(t *testing.T) {
	const (
		nilBoundValsErr = "accepted two nil bound values"
		generalBoundsErr = "only one bound value"
	)
	sm, _, err := newSortedMapFromRandRecords(1000)
	if err != nil {
		t.Fatal(err)
	}

	reversed := true
	buffSize := 64

	earlierDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	laterDate := time.Now()

	params := IterChParams{
		Reversed: reversed,
		BufSize: buffSize,
	}

	if ch, ok := sm.CustomIterCh(params); ok {
		if err := verifyRecords(ch, params.Reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestCustomIterCh failed: %v", generalBoundsErr)
	}

	params = IterChParams{
		Reversed: reversed,
		BufSize: buffSize,
		LowerBound: earlierDate,
		UpperBound: laterDate,
	}

	if ch, ok := sm.CustomIterCh(params); ok {
		if err := verifyRecords(ch, params.Reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestCustomIterCh failed: %v", generalBoundsErr)
	}

	params = IterChParams{
		Reversed: reversed,
		BufSize: buffSize,
		LowerBound: laterDate,
		UpperBound: earlierDate,
	}

	if ch, ok := sm.CustomIterCh(params); ok {
		if err := verifyRecords(ch, params.Reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestCustomIterCh failed: %v", generalBoundsErr)
	}

	reversed = false
	params = IterChParams{
		Reversed: reversed,
		BufSize: 0,
		LowerBound: laterDate,
		UpperBound: earlierDate,
	}

	if ch, ok := sm.CustomIterCh(params); ok {
		if err := verifyRecords(ch, params.Reversed); err != nil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("TestCustomIterCh failed: %v", generalBoundsErr)
	}
}

func TestIterFunc(t *testing.T) {
	sm, _, err := newSortedMapFromRandRecords(1000)
	if err != nil {
		t.Fatal(err)
	}
	sm.IterFunc(false, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestIterFunc failed: %v", nilValErr)
		}
		return true
	})
	sm.IterFunc(true, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestIterFunc failed: %v", nilValErr)
		}
		return true
	})
	i := 0
	sm.IterFunc(false, func(rec Record) bool {
		if i > 0 {
			t.Fatalf("TestIterFunc failed: %v", runawayIterErr)
		}
		i++
		return false
	})
	i = 0
	sm.IterFunc(true, func(rec Record) bool {
		if i > 0 {
			t.Fatalf("TestIterFunc failed: %v", runawayIterErr)
		}
		i++
		return false
	})
}

func TestBoundedIterFunc(t *testing.T) {
	sm, _, err := newSortedMapFromRandRecords(1000)
	if err != nil {
		t.Fatal(err)
	}

	earlierDate := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	laterDate := time.Now()

	if ok := sm.BoundedIterFunc(false, nil, nil, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(true, nil, nil, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(false, nil, laterDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(false, laterDate, nil, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(true, nil, laterDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(true, laterDate, nil, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(false, earlierDate, laterDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(true, earlierDate, laterDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(false, laterDate, earlierDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(true, laterDate, earlierDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); !ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}

	if ok := sm.BoundedIterFunc(false, laterDate, laterDate, func(rec Record) bool {
		if rec.Key == nil {
			t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
		}
		return false
	}); ok {
		t.Fatalf("TestBoundedIterFunc failed: %v", nilValErr)
	}
}