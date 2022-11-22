package gocollections

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestTimeExpiredMap(t *testing.T) {
	t.Parallel()
	var want int
	var got int
	var key = "1"
	var value = "test 1"

	tmap := NewTimeExpiredMap(1 * time.Second)
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
	})
}

func TestLoad(t *testing.T) {
	t.Parallel()
	var count = 10000

	tmap := NewTimeExpiredMap(2 * time.Second)
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

	tmap := NewTimeExpiredMap(2 * time.Second)
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

	tmap := NewTimeExpiredMap(1 * time.Second)
	defer tmap.Discard()

	tmap.AddWithDuration("1", "test 1", time.Duration(5*time.Second))

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
