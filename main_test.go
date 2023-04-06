package gocollections

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTimeExpiredList(t *testing.T) {
	t.Parallel()

	want := "value1"
	tlist := NewTimeExpiredList[string](1 * time.Second)
	defer tlist.Discard()

	tlist.Add(want)

	size := tlist.Size()
	if size != 1 {
		t.Fatalf("Expext one element in collection. But size is: %d", size)
	}

	got, err := tlist.Get(0)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("want: %s, got: %s", want, got)
	}

	time.Sleep(5 * time.Second)
	size = tlist.Size()
	if size != 0 {
		t.Fatalf("Expecting no element in collection. But got size: %d", size)
	}
}

func TestTimeExpiredList_GetAll(t *testing.T) {
	t.Parallel()

	tlist := NewTimeExpiredList[string](5 * time.Second)
	defer tlist.Discard()

	tlist.Add("value1")
	tlist.Add("value2")
	tlist.Add("value3")
	tlist.Add("value4")
	tlist.Add("value5")

	values := tlist.GetAll()
	if len(values) != 5 {
		t.Fatalf("Expected 5 values. Got: %d", len(values))
	}
}

func TestTimeExpiredList_Del(t *testing.T) {
	t.Parallel()

	tlist := NewTimeExpiredList[string](600 * time.Second)
	//defer tlist.Discard()

	tlist.Add("value1")
	tlist.Add("value2")
	tlist.Add("value3")
	tlist.Add("value4")
	tlist.Add("value5")

	// remove value3 element
	_ = tlist.Del(2)
	size := len(tlist.GetAll())
	if size != 4 {
		t.Fatalf("Expect size to be 4. But got: %d", size)
	}
	for _, v := range tlist.GetAll() {
		if v == "value3" {
			t.Fatalf("We expecte 'value3' is removed, but it was present")
		}
	}

	// remove first element value1
	_ = tlist.Del(0)
	size = len(tlist.GetAll())
	if size != 3 {
		t.Fatalf("Expect size to be 4. But got: %d", size)
	}
	for _, v := range tlist.GetAll() {
		if v == "value1" {
			t.Fatalf("We expecte 'value3' is removed, but it was present")
		}
	}

	// remove last element value5
	_ = tlist.Del(tlist.Size() - 1)
	size = len(tlist.GetAll())
	if size != 2 {
		t.Fatalf("Expect size to be 4. But got: %d", size)
	}
	for _, v := range tlist.GetAll() {
		if v == "value5" {
			t.Fatalf("We expecte 'value3' is removed, but it was present")
		}
	}

	// remove last element
	_ = tlist.Del(tlist.Size() - 1)
	size = len(tlist.GetAll())
	if size != 1 {
		t.Fatalf("Expect size to be 4. But got: %d", size)
	}

	// remove last element
	_ = tlist.Del(tlist.Size() - 1)
	size = len(tlist.GetAll())
	if size != 0 {
		t.Fatalf("Expect size to be 4. But got: %d", size)
	}

	// remove last element from empty list
	err := tlist.Del(0)
	if !errors.Is(err, ErrIndexOutOfBound) {
		t.Fatalf("Expect ErrIndexOutOfBound but got %s", err.Error())
	}
}

func TestTimeExpiredMap(t *testing.T) {
	t.Parallel()
	var want int
	var got int
	var key = "1"
	var value = "test 1"

	tmap := NewTimeExpiredMap[string, string](1 * time.Second)
	defer tmap.Discard()

	t.Run("Size", func(t *testing.T) {
		want = 0
		got = tmap.Size()
		if want != got {
			t.Errorf("want: %d, got: %d", want, got)
		}
	})

	t.Run("Add", func(t *testing.T) {
		tmap.Add(key, value)

		want = 1
		got = tmap.Size()
		if want != got {
			t.Errorf("want: %d, got: %d", want, got)
		}
	})

	t.Run("Contains", func(t *testing.T) {
		c := tmap.Contains(key)
		if !c {
			t.Errorf("key is not in the map")
		}
	})

	t.Run("Get", func(t *testing.T) {
		val, err := tmap.Get(key)
		if err != nil {
			t.Error(err)
		}
		if value != val {
			t.Errorf("want: %s, got: %s", value, val)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		time.Sleep(2 * time.Second)
		want = 0
		got = tmap.Size()
		if want != got {
			t.Errorf("want: %d, got: %d", want, got)
		}
		c := tmap.Contains(key)
		if c {
			t.Error("contains, `key is in the map, but should expire")
		}
	})
}

func TestLoad(t *testing.T) {
	t.Parallel()
	var count = 10000

	tmap := NewTimeExpiredMap[string, string](2 * time.Second)
	defer tmap.Discard()

	for i := 1; i < count+1; i++ {
		tmap.Add(strconv.Itoa(i), fmt.Sprintf("TEST %d", i))
	}

	if tmap.Size() != count {
		t.Fatalf("We expect %d number of elemets, got: %d", count, tmap.Size())
	}

	time.Sleep(4 * time.Second)

	if tmap.Size() != 0 {
		t.Fatalf("We expect all elements expired but size of map is %d", tmap.Size())
	}
}

func TestTimeExpiredMap_Del(t *testing.T) {
	t.Parallel()
	var want int
	var got int

	tmap := NewTimeExpiredMap[string, string](2 * time.Second)
	defer tmap.Discard()

	tmap.Add("1", "test 1")

	want = 1
	got = tmap.Size()
	if want != got {
		t.Errorf("want: %d, got: %d", want, got)
	}

	err := tmap.Del("1")
	if err != nil {
		t.Fatal(err)
	}
	want = 0
	got = tmap.Size()
	if want != got {
		t.Errorf("want: %d, got: %d", want, got)
	}
}

func TestTimeExpiredMap_AddWithDuration(t *testing.T) {
	t.Parallel()
	var want int
	var got int

	tmap := NewTimeExpiredMap[string, string](1 * time.Second)
	defer tmap.Discard()

	tmap.AddWithDuration("1", "test 1", 5*time.Second)

	want = 1
	got = tmap.Size()
	if want != got {
		t.Errorf("want: %d, got: %d", want, got)
	}

	time.Sleep(2 * time.Second)
	want = 1
	got = tmap.Size()
	if want != got {
		t.Errorf("want: %d, got: %d", want, got)
	}

	time.Sleep(5 * time.Second)

	want = 0
	got = tmap.Size()
	if want != got {
		t.Errorf("want: %d, got: %d", want, got)
	}
}
