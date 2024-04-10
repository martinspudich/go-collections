# go-collections

Implementation of custom go collections.

### List of collections

* TimeExpiredList
  * Elements of this list has expiration duration. After this duration elements are removed from the list.
  * When the list is created via NewTimeExpiredList function it starts goroutine which removes expired elements.
* TimeExpiredMap
  * Elements of this map has expiration duration. After this duration elements are removed from the map.
  * When the map is created via NewTimeExpiredMap function it starts goroutine which removes expired elements.

### TimeExpiredMap

### Time ExpiredList

basic usage:
```go
import (
	goc "github.com/martinspudich/go-collections"
)

func main() {

    tlist := goc.NewTimeExpiredList(2 * time.Second) // creates new TimeExpiredMap, it starts new goroutine
    defer tlist.Discard()                           // stops goroutine and discards internal data map

    tlist.Add("test 1") // adds element with value "test 1"
}
```

basic usage:
```go
import (
	goc "github.com/martinspudich/go-collections"
)

func main() {

    tmap := goc.NewTimeExpiredMap(2 * time.Second) // creates new TimeExpiredMap, it starts new goroutine
    defer tmap.Discard()                           // stops goroutine and discards internal data map

    tmap.Add("1", "test 1") // adds element with key "1" and value "test 1"
}
```

#### Using expired elemenet channel

This channel is used only if we configure `expiredElChanSize` bigger the 0. Default is 0. By default, no expired elements are sent
to channel `ExpiredElChan`.

```go
import (
	goc "github.com/martinspudich/go-collections"
)

func main() {
    // creates new TimeExpiredMap, it starts new goroutine
    tlist := goc.NewTimeExpiredList(100 * time.Millisecond, Config{
        CleanJobInterval:  200 * time.Millisecond, // default 60s.
        ExpiredElChanSize: 100, // if it's bigger than 0, than expired element channel is used
    })
    defer tlist.Discard() // stops goroutine and discards internal data map

    tlist.Add("test 1") // adds element with value "test 1"
	
	time.Sleep(300 * time.Millisecond)
	<- tlist.ExpiredElChan() // Receive expired element.
}
```

### RELEASE NOTES

#### 0.4.0
* Add expired element channel
  * to this channel we add expired elements, when elements expired
  * so it's possible to react on the event