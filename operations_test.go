package odatas

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/couchbase/gocb"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	if _, _, err := testInsert(); err != nil {
		t.Fatal(err)
	}
}

func TestWritePtrValue(t *testing.T) {
	ws := generate()
	_, err := th.Insert(context.Background(), "webshop", &ws)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWritePrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	_, err := th.Insert(context.Background(), "webshop", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

func TestWritePrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	_, err := th.Insert(context.Background(), "webshop", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

func TestWriteNotExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	_, err := th.Insert(context.Background(), "member", s)
	if err != nil {
		t.Error("Missing error")
	}
	log.Println(err)
}

//func TestWriteExpectError(t *testing.T) {
//	s := struct {
//		name string
//	}{name: "Jackson"}
//	id, err := th.Insert(context.Background(), "member", s)
//	if err != nil {
//		t.Error("Missing error")
//	}
//
//	q:= func(typ, id string, ptr interface{}, ttl int) (gocb.Cas, error) {
//		documentID := typ + "::" + id
//		return h.state.bucket.Insert(documentID, ptr, 0), nil
//	}
//	_, errDuplicateInsert := th.write(context.Background(), "member", id, s,q)
//	if errDuplicateInsert == nil {
//		t.Error("error missing", errDuplicateInsert)
//	}
//}

func testInsert() (webshop, string, error) {
	ws := generate()
	id, err := th.Insert(context.Background(), "webshop", ws)
	return ws, id, err
}

func TestRead(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	if err := th.Get(context.Background(), "webshop", id, &ws); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", ws)
}

func TestReadPrimitivePtrNil(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Insert(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{Job: nil}
	errGet := th.Get(context.Background(), "webshop", id, &testGet)
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
	id, errInsert := th.Insert(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	b := "b"
	var testGet = wtyp{Job: &b}
	errGet := th.Get(context.Background(), "webshop", id, &testGet)
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
	id, errInsert := th.Insert(context.Background(), "webshop", test)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
	errGet := th.Get(context.Background(), "webshop", id, &testGet)
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
	id, errInsert := th.Insert(context.Background(), "webshop", testInsert)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
	errGet := th.Get(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("error")
	}
	assert.NotEqual(t, testInsert, testGet, "They should be not equal")

}

func TestIDNotFoundError(t *testing.T) {
	id := "123"
	ws := webshop{}
	if err := th.Get(context.Background(), "webshop", id, &ws); err == nil {
		t.Error("read with invalid ID")
	}
}

func TestPingNilService(t *testing.T) {
	pingReport, err := th.Ping(context.Background(), nil)
	if err != nil {
		t.Error("error", err)
	}
	fmt.Printf("%+v\n", *pingReport)
}

func TestPingAllService(t *testing.T) {
	services := make([]gocb.ServiceType, 5)
	services = append(services, gocb.MemdService)

	pingReport, err := th.Ping(context.Background(), []gocb.ServiceType{gocb.MemdService, gocb.MgmtService, gocb.CapiService, gocb.N1qlService, gocb.FtsService, gocb.CbasService})
	if err != nil {
		t.Error("error", err)
	}
	fmt.Printf("%+v\n", pingReport)
}

func TestTouch(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{
		Email: "",
		Product: &product{
			ID:          "",
			UserID:      "",
			StoreID:     "",
			Name:        "",
			Description: "",
			Slug:        "",
			Price:       0,
			SalePrice:   0,
			CurrencyID:  0,
			OnSale:      0,
			Status:      "",
		},
		Store: &store{
			ID:          "",
			UserID:      "",
			Name:        "",
			Description: "",
		},
	}
	if err := th.Touch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
}

func TestGetAndTouch(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{
		Email: "",
		Product: &product{
			ID:          "",
			UserID:      "",
			StoreID:     "",
			Name:        "",
			Description: "",
			Slug:        "",
			Price:       0,
			SalePrice:   0,
			CurrencyID:  0,
			OnSale:      0,
			Status:      "",
		},
		Store: &store{
			ID:          "",
			UserID:      "",
			Name:        "",
			Description: "",
		},
	}
	if err := th.GetAndTouch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
}

//func TestUpsert(t *testing.T) {
//	ws, ID, err := testInsert()
//	if err != nil {
//		t.Fatal(err)
//	}
//	updateableWs := *(&ws)
//	updateableWs.Email = gofakeit.Email()
//	updateableWs.Product.Name = gofakeit.Name()
//
//	if err := th.Upsert(ID, "webshop", updateableWs, 0); err != nil {
//		t.Fatal(err)
//	}
//	if ws.Email == updateableWs.Email {
//		t.Error("Update error at Email")
//	}
//	if ws.Product.Name == updateableWs.Product.Name {
//		t.Error("Update error at Product's Name")
//	}
//
//}

func TestRemove(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	if err := th.Remove(context.Background(), "webshop", ID, &webshop{}); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkInsertEmb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = testInsert()
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = th.Insert(context.Background(), "webshop", generate())
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := th.Insert(context.Background(), "webshop", generate())
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Get(context.Background(), "webshop", ID, webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkGetEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, id, _ := testInsert()
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Get(context.Background(), "webshop", id, &webshop{})
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
		id, _ := th.Insert(context.Background(), "job", job)
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		var jobRead jobtyp
		start := time.Now()
		_ = th.Get(context.Background(), "job", id, &jobRead)
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
