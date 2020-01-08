package bucket

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/rs/xid"
	"github.com/volatiletech/null"

	"github.com/couchbase/gocb"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	if _, id, err := testInsert(); err != nil || id == "" {
		t.Fatal(err)
	}
}

func TestInsertNullString(t *testing.T) {
	type nullStrTest struct {
		Bin         int         `json:"bin"`
		CardBrand   string      `json:"card_brand"`
		IssuingBank string      `json:"issuing_bank"`
		CardType    null.String `json:"card_type"`
	}
	cardTypeWrite := nullStrTest{
		Bin:         50003,
		CardBrand:   "VISA",
		IssuingBank: "",
		CardType:    null.String{String: "US", Valid: true},
	}
	_, _, err := th.Insert(context.Background(), "card_type", "", cardTypeWrite, 0)
	if err != nil {
		t.Error(err)
	}

}

func TestInsertCustomID(t *testing.T) {
	cID := xid.New().String() + "Faswwq123942390**12312_+"
	ws := generate()
	_, id, err := th.Insert(context.Background(), "webshop", cID, &ws, 0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, cID, id, "should be equal")
}

func TestInsertDocumentTypeNotFoundState(t *testing.T) {
	_ = th.state.deleteType("webshop")
	if _, _, err := testInsert(); err != nil {
		t.Error(err)
	}
}

func testInsert() (webshop, string, error) {
	ws := generate()
	_, id, err := th.Insert(context.Background(), "webshop", "", ws, 0)
	return ws, id, err
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
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	if err := th.Touch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
}

func TestTouchNonPointerInputExpectError(t *testing.T) {
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	if err := th.Touch(context.Background(), "webshop", ID, ws, 10); err != ErrInvalidGetDocumentTypesParam {
		t.Errorf("error should be %s instead of %s", ErrInvalidGetDocumentTypesParam, err)
	}
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

func TestUpsertPtrValueSameID(t *testing.T) {
	ws := generate()
	id := xid.New().String()
	_, _, err := th.Upsert(context.Background(), "webshop", id, &ws, 0)
	if err != nil {
		t.Fatal(err)
	}
	wsD := generate()
	_, _, errD := th.Upsert(context.Background(), "webshop", id, &wsD, 1)
	if errD != nil {
		t.Fatal(errD)
	}
	errGet := th.Get(context.Background(), "webshop", id, &ws)
	if errGet != nil {
		t.Fatal(errGet)
	}
	assert.Equal(t, wsD, ws, "should be equal")
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

func TestUpsertTypeNotFoundExpectError(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	ws := generate()
	_ = th.state.deleteType("product")
	_, _, err = th.Upsert(context.Background(), "webshop", id, ws, 0)
	if err != nil {
		t.Error(err)
	}

}

func testUpsert(id string) (webshop, string, error) {
	ws := generate()
	_, id, err := th.Upsert(context.Background(), "webshop", id, ws, 0)
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

func TestRemoveInvalidInput(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	if err := th.Remove(context.Background(), "webshop", ID, webshop{}); err != ErrInvalidGetDocumentTypesParam {
		t.Errorf("error should be %s instead of %s", ErrInvalidGetDocumentTypesParam, err)
	}
}

func TestRemoveDocumentDoesntExist(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	if err = th.Remove(context.Background(), "webshop", id, &webshop{}); err != nil {
		t.Fatal(err)
	}
	if err = th.Remove(context.Background(), "webshop", id, &webshop{}); err != gocb.ErrKeyNotFound {
		t.Errorf("error should be %s instead of %s", gocb.ErrKeyNotFound, err)

	}
}

func BenchmarkInsertEmb(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = testInsert()
	}
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = th.Insert(context.Background(), "webshop", "", generate(), 0)
	}
}

func BenchmarkRemoveEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		_, ID, _ := testInsert()
		split := strings.Split(ID, "::")
		b.StartTimer()
		_ = th.Remove(context.Background(), split[1], split[0], &webshop{})
	}
}
