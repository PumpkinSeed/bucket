package bucket

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestInsertNilEmbeddedStruct(t *testing.T) {
	ws := generate()
	ws.Product = nil
	cas, id, err := th.Insert(context.Background(), "webshop", "", ws, 0)
	if err != nil || id == "" {
		t.Error(err)
	}
	if len(cas) != 2 {
		t.Errorf("Cas should store 2 elements, instead of %d", len(cas))
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

func TestInsertPtrValue(t *testing.T) {
	ws := generate()
	cas, id, err := th.Insert(context.Background(), "webshop", "", &ws, 0)
	if err != nil || id == "" {
		t.Fatal(err)
	}
	if len(cas) != 4 {
		t.Errorf("Cas should store 4 elements, instead of %d", len(cas))
	}
}

func TestInsertPrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	cas, id, err := th.Insert(context.Background(), "webshop", "", s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))
	}
}

func TestInsertPrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	cas, id, err := th.Insert(context.Background(), "webshop", "", s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))
	}
}

func TestInsertNonExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	cas, id, err := th.Insert(context.Background(), "member", "", s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))
	}
}

func TestInsertExpectDuplicateError(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	ctx := context.Background()
	id := xid.New().String()
	_, _, err := th.Insert(ctx, "member", id, s, 0)
	if err != nil {
		t.Error("Missing error")
	}
	_, _, errDuplicateInsert := th.Insert(ctx, "member", id, s, 0)
	if errDuplicateInsert == nil {
		t.Error("error missing", errDuplicateInsert)
	}
	assert.EqualValues(t, "key already exists, if a cas was provided the key exists with a different cas", errDuplicateInsert.Error(), "wrong error msg")
}

func TestInsertEmptyRefTag(t *testing.T) {
	ws := generate()
	s := struct {
		Name    string   `json:"name"`
		Product *product `json:"product" cb_referenced:""`
	}{
		Name:    "Missing",
		Product: ws.Product,
	}
	_, _, err := th.Insert(context.Background(), "name", "", s, 0)
	if err != ErrEmptyRefTag {
		t.Errorf("Error should be %s instead of %s", ErrEmptyRefTag, err)
	}
}

func TestInsertEmbeddedStructExpectKeyAlreadyExistError(t *testing.T) {
	prod := product{}
	ctx := context.Background()
	_, id, err := th.Insert(ctx, "product", "", prod, 0)
	if err != nil {
		t.Fatal(err)
	}
	ws := generate()
	if _, _, err := th.Insert(ctx, "webshop", id, ws, 0); err != gocb.ErrKeyExists {
		t.Errorf("error should be %s instead of %s", gocb.ErrKeyExists, err)

	}
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

func TestGetNullString(t *testing.T) {
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
	_, id, err := th.Insert(context.Background(), "card_type", "", cardTypeWrite, 0)
	if err != nil {
		t.Error(err)
	}
	cardTypeRead := nullStrTest{}
	if err := th.Get(context.Background(), "card_type", id, &cardTypeRead); err != nil {
		t.Error(err)
	}
	assert.Equal(t, cardTypeWrite, cardTypeRead, "should be equal")
}

func TestGetNilEmbeddedStruct(t *testing.T) {
	wsInsert := generate()
	wsInsert.Product = nil
	ctx := context.Background()
	typ := "webshop"
	id := ""
	_, id, err := th.Insert(ctx, typ, id, wsInsert, 0)
	if err != nil {
		t.Error(err)
	}
	wsGet := &webshop{}
	errGet := th.Get(ctx, typ, id, wsGet)
	if errGet != nil {
		t.Error(errGet)
	}
	assert.Equal(t, &wsInsert, wsGet, "should be equal")
}

func TestGetPrimitivePtrNil(t *testing.T) {
	a := "a"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	_, id, errInsert := th.Insert(context.Background(), "webshop", "", test, 0)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
	errGet := th.Get(context.Background(), "webshop", id, &testGet)
	if errGet != nil {
		t.Error("Error")
	}
	assert.Equal(t, test, testGet, "They should be equal")
}

func TestGetPrimitivePtr(t *testing.T) {
	a := "artist"
	type wtyp struct {
		Job *string `json:"name,omitempty"`
	}
	test := wtyp{Job: &a}
	_, id, errInsert := th.Insert(context.Background(), "webshop", "", test, 0)
	if errInsert != nil {
		t.Error("Error")
	}
	var testGet = wtyp{}
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
	_, id, errInsert := th.Insert(context.Background(), "webshop", "", test, 0)
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
	_, id, errInsert := th.Insert(context.Background(), "webshop", "", testInsert, 0)
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
	ws := webshop{}
	if err := th.Get(context.Background(), "webshop", "123", &ws); err == nil {
		t.Error("read with invalid ID")
	}
}

func TestGetInvalidInput(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	wsGet := webshop{}
	if err := th.Get(context.Background(), "webshop", id, wsGet); err != ErrInputStructPointer {
		t.Errorf("error should be %s instead of %s", ErrInputStructPointer, err)
	}
}

func TestGetTypeNotFoundExpectError(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	_ = th.state.deleteType("webshop")
	if err := th.Get(context.Background(), "webshop", id, webshop{}); err != ErrDocumentTypeDoesntExists {
		t.Errorf("error should be %s instead of %s", ErrDocumentTypeDoesntExists, err)
	}
}

func TestGetEmptyRefTagExpectErr(t *testing.T) {
	type wsInsert struct {
		Token   string   `json:"token"`
		Product *product `json:"product" cb_referenced:"product"`
	}
	websh := wsInsert{
		Token: "",
		Product: &product{
			Name:        "testprod",
			Description: "description",
			Price:       1221,
			CurrencyID:  923,
		},
	}
	type wsGet struct {
		Token   string   `json:"token"`
		Product *product `json:"product" cb_referenced:""`
	}
	ctx := context.Background()
	_, id, err := th.Insert(ctx, "webshop", "", websh, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := th.Get(ctx, "webshop", id, &wsGet{}); err != ErrEmptyRefTag {
		t.Errorf("error should be %s instead of %s", ErrEmptyRefTag, err)
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
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	if err := th.Touch(context.Background(), "webshop", ID, &ws, 10); err != nil {
		t.Error("error", err)
	}
}

func TestTouchDocumentTypeNotFoundExpectError(t *testing.T) {
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	_ = th.state.deleteType("product")

	if err := th.Touch(context.Background(), "webshop", ID, &ws, 10); err != ErrDocumentTypeDoesntExists {
		t.Errorf("error should be %s instead of %s", ErrDocumentTypeDoesntExists, err)
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
func TestGetAndTouchDocumentTypeNotFoundExpectError(t *testing.T) {
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	_ = th.state.deleteType("product")

	if err := th.GetAndTouch(context.Background(), "webshop", ID, &ws, 10); err != ErrDocumentTypeDoesntExists {
		t.Errorf("error should be %s instead of %s", ErrDocumentTypeDoesntExists, err)
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

func TestUpsertPtrValueNewID(t *testing.T) {
	ws := generate()
	cas, id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), &ws, 0)
	if err != nil || id == "" {
		t.Fatal(err)
	}
	if len(cas) != 4 {
		t.Errorf("Cas should store 4 element, instead of %d", len(cas))

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

func TestUpsertPrimitivePtr(t *testing.T) {
	asd := "asd"
	s := struct {
		Name *string `json:"name,omitempty"`
	}{Name: &asd}
	cas, id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))

	}
}

func TestUpsertPrimitivePtrNil(t *testing.T) {
	s := struct {
		Name *string `json:"name,omitempty"`
	}{}
	cas, id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))
	}
}

func TestUpsertNonExportedField(t *testing.T) {
	s := struct {
		name string
	}{name: "Jackson"}
	cas, id, err := th.Upsert(context.Background(), "webshop", xid.New().String(), s, 0)
	if err != nil || id == "" {
		t.Error("Missing error")
	}
	if len(cas) != 1 {
		t.Errorf("Cas should store 1 element, instead of %d", len(cas))
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

func TestRemoveDocumentKeyNotFoundExpectError(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	err = th.state.deleteType("product")
	if err != nil {
		t.Fatal(err)
	}
	err = th.Remove(context.Background(), "webshop", id, &webshop{})
	if err != ErrDocumentTypeDoesntExists {
		t.Errorf("error should be %s instead of %s", ErrDocumentTypeDoesntExists, err)
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

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, ID, _ := th.Insert(context.Background(), "webshop", "", generate(), 0)
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Get(context.Background(), "webshop", ID, webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkGetEmbedded(b *testing.B) {
	b.StopTimer()
	_, id, _ := testInsert()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = th.Get(context.Background(), "webshop", id, &webshop{})
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
		_, id, _ := th.Insert(context.Background(), "job", "", job, 0)
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
