package thumb

import "strconv"

// MaxAge represents a cache TTL in seconds.
type MaxAge int

// String returns the cache TTL in seconds as string.
func (a MaxAge) String() string {
	return strconv.Itoa(int(a))
}

var (
	CacheMaxAge MaxAge = 2592000
	CachePublic        = false
)
