package examples

import (
	"math/rand"
	"time"
)

func RandomReqNumber() int64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(r.Intn(9999999) + 1)
}

// Ptr returns a pointer to v helper for optional struct fields.
func Ptr[T any](v T) *T { return &v }
