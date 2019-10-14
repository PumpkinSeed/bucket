package bucket

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_EffGet(t *testing.T) {
	wsInsert, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	wsGet := webshop{}
	if err := th.EffGet(context.Background(), "webshop", id, &wsGet); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, wsInsert, wsGet, "should be equal")
}

func BenchmarkHandler_EffGet(b *testing.B) {
	b.StopTimer()
	_, id, err := testInsert()
	if err != nil {
		b.Fatal(err)
	}

	wsGet := webshop{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := th.EffGet(context.Background(), "webshop", id, &wsGet); err != nil {
			b.Fatal(err)
		}
	}
}
