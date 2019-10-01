package bucket

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	if _, id, err := testInsert(); err != nil || id == "" {
		t.Fatal(err)
	}
}

func TestInsertCustomID(t *testing.T) {
	cID := xid.New().String() + "Faswwq123942390**12312_+"
	ws := generate()
	id, err := th.Insert(context.Background(), "webshop", cID, &ws)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, cID, id, "should be equal")
}

func TestInsertPtrValue(t *testing.T) {
	ws := generate()
	id, err := th.Insert(context.Background(), "webshop", "", &ws)
	if err != nil || id == "" {
		t.Fatal(err)
	}
}

func TestInsertPrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	id, err := th.Insert(context.Background(), "webshop", "", s)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestInsertPrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	id, err := th.Insert(context.Background(), "webshop", "", s)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestInsertNonExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	id, err := th.Insert(context.Background(), "member", "", s)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestInsertExpectDuplicateError(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	ctx := context.Background()
	id := xid.New().String()
	_, err := th.Insert(ctx, "member", id, s)
	if err != nil {
		t.Error("Missing error")
	}
	_, errDuplicateInsert := th.Insert(ctx, "member", id, s)
	if errDuplicateInsert == nil {
		t.Error("error missing", errDuplicateInsert)
	}
	assert.EqualValues(t, "key already exists, if a cas was provided the key exists with a different cas", errDuplicateInsert.Error(), "wrong error msg")
}

func testInsert() (webshop, string, error) {
	ws := generate()
	id, err := th.Insert(context.Background(), "webshop", "", ws)
	return ws, id, err
}

func TestGet(t *testing.T) {
	wsInsert, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	wsGet := webshop{}
	if err := th.Get(context.Background(), "webshop", id, &wsGet); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, wsInsert, wsGet, "should be equal")
}

func TestGetPrimitivePtrNil(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Insert(context.Background(), "webshop", "", test)
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

func TestGetPrimitivePtr(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Insert(context.Background(), "webshop", "", test)
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

func TestGetNonPointerInput(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	id, errInsert := th.Insert(context.Background(), "webshop", "", test)
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

func TestGetNonExportedField(t *testing.T) {
	a := "helder"
	type wtyp struct {
		job string
	}
	testInsert := wtyp{job: a}
	id, errInsert := th.Insert(context.Background(), "webshop", "", testInsert)
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

func TestGetIDNotFoundError(t *testing.T) {
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
	for _, service := range pingReport.Services {
		assert.Equal(t, service.Success, true, "should be true")
		fmt.Printf("%+v\n", service)

	}
}

func TestPingAllService(t *testing.T) {
	services := make([]gocb.ServiceType, 5)
	services = append(services, gocb.MemdService)

	pingReport, err := th.Ping(context.Background(), []gocb.ServiceType{gocb.MemdService, gocb.MgmtService, gocb.CapiService, gocb.N1qlService, gocb.FtsService, gocb.CbasService})
	if err != nil {
		t.Error("error", err)
	}
	for _, service := range pingReport.Services {
		if service.Service == gocb.CbasService {
			assert.Equal(t, service.Success, false, "should be false,service is missing")
		} else {
			assert.Equal(t, service.Success, true, "should be true")

		}
		fmt.Printf("%+v\n", service)
	}
}

func TestTouch(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	if err := th.Touch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
}

func TestGetAndTouch(t *testing.T) {
	webshopInsert, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	if err := th.GetAndTouch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
	assert.Equal(t, webshopInsert, ws, "should be equal")
}

func TestUpsertNewID(t *testing.T) {
	if _, id, err := testUpsert(xid.New().String()); err != nil || id == "" {
		t.Fatal(err)
	}
}

func TestUpsertSameID(t *testing.T) {
	id := xid.New().String()
	if _, _, err := testUpsert(id); err != nil {
		t.Fatal(err)
	}
	if _, _, err := testUpsert(id); err != nil {
		t.Fatal(err)
	}
}

func TestUpsertPtrValueNewID(t *testing.T) {
	ws := generate()
	id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), &ws, 0)
	if err != nil || id == "" {
		t.Fatal(err)
	}
}

func TestUpsertPtrValueSameID(t *testing.T) {
	ws := generate()
	id := xid.New().String()
	_, err := th.Upsert(context.Background(), "webshop", id, &ws, 0)
	if err != nil {
		t.Fatal(err)
	}
	wsD := generate()
	_, errD := th.Upsert(context.Background(), "webshop", id, &wsD, 1)
	if errD != nil {
		t.Fatal(errD)
	}
	errGet := th.Get(context.Background(), "webshop", id, &ws)
	if errGet != nil {
		t.Fatal(errGet)
	}
	assert.Equal(t, wsD, ws, "should be equal")
}

func TestUpsertPrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestUpsertPrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestUpsertNonExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
}

func TestUpsertEmptyID(t *testing.T) {
	_, id, err := testUpsert("")
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Error("invalid id")
	}
}

func testUpsert(id string) (webshop, string, error) {
	ws := generate()
	id, err := th.Upsert(context.Background(), "webshop", id, ws, 0)
	return ws, id, err
}

func TestRemove(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	if err := th.Remove(context.Background(), "webshop", ID, &webshop{}); err != nil {
		t.Fatal(err)
	}
	if err := th.Get(context.Background(), "webshop", ID, &webshop{}); err != nil {
		assert.Equal(t, gocb.ErrKeyNotFound, err, "error")
	}
}

func BenchmarkInsertEmb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = testInsert()
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = th.Insert(context.Background(), "webshop", "", generate())
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := th.Insert(context.Background(), "webshop", "", generate())
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
		id, _ := th.Insert(context.Background(), "job", "", job)
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
