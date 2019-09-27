package odatas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestWriteNotExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	_, err := th.Write(context.Background(), "member", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

func TestWriteExpectError(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	id, err := th.Write(context.Background(), "member", s)
	if err != nil {
		t.Error("Missing error")
	}
	_, errDuplicateInsert := th.write(context.Background(), "member", id, s)
	if errDuplicateInsert == nil {
		t.Error("error missing", errDuplicateInsert)
	}
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

func TestIDNotFoundError(t *testing.T) {
	id := "123"
	ws := webshop{}
	if err := th.Read(context.Background(), "webshop", id, &ws); err == nil {
		t.Fail()
	}
}

func TestReadPrimitivePtrNil(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Write(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{Job: nil}
	errGet := th.Read(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("Error")
	}
	assert.Equal(t, test, testGet, "They should be equal")
}

func TestReadPrimitivePtr(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Write(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	b := "b"
	var testGet = wtyp{Job: &b}
	errGet := th.Read(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("Error")
	}
	assert.Equal(t, test, testGet, "They should be equal")
}

func TestReadNonPointerInput(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Write(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
	errGet := th.Read(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("error")
	}
	assert.Equal(t, test, testGet, "They should be equal")
}

func TestReadNotExportedField(t *testing.T) {
	a := "helder"
	type wtyp struct {
		job string
	}
	testInsert := wtyp{job: a}
	id, errInsert := th.Write(context.Background(), "webshop", testInsert)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
	errGet := th.Read(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("error")
	}
	assert.NotEqual(t, testInsert, testGet, "They should be not equal")

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

func BenchmarkGetPtr(b *testing.B) {
	type jobtyp struct {
		Job *string `json:"job,omitempty"`
	}
	j := "helder"
	for i := 0; i < b.N; i++ {
		job := jobtyp{Job: &j}
		startInsert := time.Now()
		id, _ := th.Write(context.Background(), "job", job)
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		var jobRead jobtyp
		start := time.Now()
		_ = th.Read(context.Background(), "job", id, &jobRead)
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
