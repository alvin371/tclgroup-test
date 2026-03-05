package port

import "time"

// Clock abstracts time retrieval for testability.
type Clock interface {
	Now() time.Time
}
