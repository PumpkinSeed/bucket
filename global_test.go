package odatas

import (
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
)

var h Handler

func init() {
	gofakeit.Seed(time.Now().UnixNano())

	h = defaultHandler()

	start := time.Now()
	if err := h.GetManager().Flush(); err != nil {
		fmt.Printf("Turn on flush in bucket: %+v\n", err)
	}
	fmt.Printf("Bucket flushed: %v\n", time.Since(start))

	for j := 0; j < 10000; j++ {
		instance := newTestStruct1()
		_, _ = h.bucket.Insert(instance.Token, instance, 0)
	}
	fmt.Printf("Connection setup, data seeded %v\n", time.Since(start))
}
