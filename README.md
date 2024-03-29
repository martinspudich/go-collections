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

    tmap := goc.NewTimeExpiredList(2 * time.Second) // creates new TimeExpiredMap, it starts new goroutine
    defer tmap.Discard()                           // stops goroutine and discards internal data map

    tmap.Add("test 1") // adds element with value "test 1"
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