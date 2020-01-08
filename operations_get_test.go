package bucket

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null"
)

func TestHandler_Get(t *testing.T) {
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

func TestHandler_GetNullString(t *testing.T) {
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

func TestHandler_GetNilEmbeddedStruct(t *testing.T) {
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

func TestHandler_GetPrimitivePtrNil(t *testing.T) {
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

func TestHandler_GetPrimitivePtr(t *testing.T) {
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

func TestHandler_GetNonPointerInput(t *testing.T) {
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

func TestHandler_GetNonExportedField(t *testing.T) {
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

func TestHandler_GetIDNotFoundError(t *testing.T) {
	ws := webshop{}
	if err := th.Get(context.Background(), "webshop", "123", &ws); err == nil {
		t.Error("read with invalid ID")
	}
}

func TestHandler_GetInvalidInput(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	wsGet := webshop{}
	if err := th.Get(context.Background(), "webshop", id, wsGet); err != ErrInputStructPointer {
		t.Errorf("error should be %s instead of %s", ErrInputStructPointer, err)
	}
}

func TestHandler_GetEmptyRefTagExpectErr(t *testing.T) {
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

func TestHandler_GetAndTouch(t *testing.T) {
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

func BenchmarkHandler_Get(b *testing.B) {
	b.StopTimer()
	_, id, err := testInsert()
	if err != nil {
		b.Fatal(err)
	}

	wsGet := webshop{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if err := th.Get(context.Background(), "webshop", id, &wsGet); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHandler_GetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, ID, _ := th.Insert(context.Background(), "webshop", "", generate(), 0)
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Get(context.Background(), "webshop", ID, webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkHandler_GetEmbedded(b *testing.B) {
	b.StopTimer()
	_, id, _ := testInsert()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = th.Get(context.Background(), "webshop", id, &webshop{})
	}
}

func BenchmarkHandler_GetPtr(b *testing.B) {
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
