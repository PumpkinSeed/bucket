package odatas

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
)

func Test(t *testing.T) {
	if _, _, err := testInsert(); err != nil {
		t.Fatal(err)
	}
}

func testInsert() (webshop, string, error) {
	ws := generate()
	ID, err := th.Write(ws, "webshop")
	return ws, ID, err
}

func TestRead(t *testing.T) {
	_, id, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{}
	//splitedID := strings.Split(ID, "::")
	if err := th.Read("webshop", id, &ws); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", ws)
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
	splitedID := strings.Split(ID, "::")
	if err := th.Touch(splitedID[1], splitedID[0], &ws, 10); err != nil {
		t.Fail()
	}
	fmt.Printf("%+v\n", ws)
}

func TestUpsert(t *testing.T) {
	ws, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	updateableWs := *(&ws)
	updateableWs.Email = gofakeit.Email()
	updateableWs.Product.Name = gofakeit.Name()

	if err := th.Upsert(ID, "webshop", updateableWs, 0); err != nil {
		t.Fatal(err)
	}
	if ws.Email == updateableWs.Email {
		t.Error("Update error at Email")
	}
	if ws.Product.Name == updateableWs.Product.Name {
		t.Error("Update error at Product's Name")
	}

}

func TestRemove(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}
	split := strings.Split(ID, "::")
	if err := th.Remove(split[1], split[0], &webshop{}); err != nil {
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
		_, _ = th.Write(generate(), "webshop")
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := th.Write(generate(), "webshop")
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = th.Read(ID, "webshop", webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}

func BenchmarkGetEmbedded(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		_, ID, _ := testInsert()
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		split := strings.Split(ID, "::")
		start := time.Now()
		_ = th.Read(split[1], split[0], &webshop{})
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
		_ = th.Remove(split[1], split[0], &webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}
