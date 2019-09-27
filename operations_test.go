package odatas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	if _, _, err := testInsert(); err != nil {
		t.Fatal(err)
	}
}

func TestWritePtrValue(t *testing.T) {
	ws := generate()
	_, err := th.Write(context.Background(), "webshop", &ws)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWritePrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	_, err := th.Write(context.Background(), "webshop", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

func TestWritePrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	_, err := th.Write(context.Background(), "webshop", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

func testInsert() (webshop, string, error) {
	ws := generate()
	id, err := th.Write(context.Background(), "webshop", ws)
	return ws, id, err
}

func TestRead(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	if err := th.Read(context.Background(), "webshop", id, &ws); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", ws)
}

func TestReadPrimitivePtrNil(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	w := wtyp{Job: &a}
	id, errInsert := th.Write(context.Background(), "webshop", w)
	if errInsert != nil {
		t.Error("Error")
	}
	var ww = wtyp{Job: nil}
	errGet := th.Read(context.Background(), "webshop", id, &ww)
	if errGet != nil {
		t.Error("Error")
	}
	fmt.Println(*ww.Job)
}

func TestReadPrimitivePtr(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	w := wtyp{Job: &a}
	id, errInsert := th.Write(context.Background(), "webshop", w)
	if errInsert != nil {
		t.Error("Error")
	}
	b := "b"
	var ww = wtyp{Job: &b}
	errGet := th.Read(context.Background(), "webshop", id, &ww)
	if errGet != nil {
		t.Error("Error")
	}
	fmt.Println(*ww.Job)
}

func BenchmarkInsertEmb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = testInsert()
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = th.Write(context.Background(), "webshop", generate())
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := th.Write(context.Background(), "webshop", generate())
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Read(context.Background(), "webshop", ID, webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkGetEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, id, _ := testInsert()
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Read(context.Background(), "webshop", id, &webshop{})
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
		_ = th.Remove(context.Background(), split[1], split[0], &webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}
