package odatas

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	if _, _, err := testInsert(); err != nil {
		t.Fatal(err)
	}
}

func testInsert() (webshop, string, error) {
	ws := generate()
	id, err := th.Write(context.Background(), "webshop",ws)
	return ws, id, err
}

func TestRead(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	if err := th.Read(context.Background(),"webshop", id, &ws); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", ws)
}

func BenchmarkInsertEmb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = testInsert()
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = th.Write(context.Background(), "webshop",generate())
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := th.Write(context.Background(), "webshop",generate())
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Read(context.Background(),"webshop", ID, webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkGetEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, id, _ := testInsert()
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Read(context.Background(),"webshop", id, &webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkRemoveEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, ID, _ := testInsert()
		fmt.Printf("Insert: %vns\tRemove: ", time.Since(startInsert).Nanoseconds())
		split := strings.Split(ID, "::")
		start := time.Now()
		_ = th.Remove(context.Background(),split[1], split[0], &webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}
