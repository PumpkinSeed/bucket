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
	ws := webshop{
		ID:       23,
		RoleID:   25,
		Name:     "Test",
		Email:    "test@test.com",
		Password: "asd",
		Product: product{
			ID:          34,
			UserID:      44,
			StoreID:     55,
			Name:        "laptop",
			Description: "its a laptop",
			Slug:        "laptop",
			Price:       1233,
			SalePrice:   1400,
			CurrencyID:  2,
			OnSale:      123,
			Status:      "active",
		},
		Store: store{
			ID:          55,
			UserID:      44,
			Name:        "productshop",
			Description: "Product shop",
		},
	}
	ID, err := h.Write(ws, "webshop")
	return ws, ID, err
}

type webshop struct {
	ID       int     `json:"id"`
	RoleID   int     `json:"role_id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Product  product `json:"product"`
	Store    store   `json:"store,omitempty"`
}

type product struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	StoreID     int    `json:"store_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Price       int64  `json:"price"`
	SalePrice   int64  `json:"sale_price"`
	CurrencyID  int    `json:"currency_id"`
	OnSale      int    `json:"on_sale"`
	Status      string `json:"status"`
}

type store struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func TestRead(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{
		ID:       0,
		RoleID:   0,
		Name:     "",
		Email:    "",
		Password: "",
		Product: product{
			ID:          0,
			UserID:      0,
			StoreID:     0,
			Name:        "",
			Description: "",
			Slug:        "",
			Price:       0,
			SalePrice:   0,
			CurrencyID:  0,
			OnSale:      0,
			Status:      "",
		},
		Store: store{
			ID:          0,
			UserID:      0,
			Name:        "",
			Description: "",
		},
	}
	splitedID := strings.Split(ID, "::")
	if err := h.Read(splitedID[1], splitedID[0], &ws); err != nil {
		t.Fail()
	}
	fmt.Printf("%+v\n", ws)
}

func TestTouch(t *testing.T) {
	_, ID, err := testInsert()
	if err != nil {
		t.Fatal(err)
	}

	ws := webshop{
		ID:       0,
		RoleID:   0,
		Name:     "",
		Email:    "",
		Password: "",
		Product: product{
			ID:          0,
			UserID:      0,
			StoreID:     0,
			Name:        "",
			Description: "",
			Slug:        "",
			Price:       0,
			SalePrice:   0,
			CurrencyID:  0,
			OnSale:      0,
			Status:      "",
		},
		Store: store{
			ID:          0,
			UserID:      0,
			Name:        "",
			Description: "",
		},
	}
	splitedID := strings.Split(ID, "::")
	if err := h.Touch(splitedID[1], splitedID[0], &ws, 10); err != nil {
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
	updateableWs.Name = gofakeit.Name()
	updateableWs.Email = gofakeit.Email()
	updateableWs.Product.Name = gofakeit.Name()

	if err := h.Upsert(ID, "webshop", updateableWs, 0); err != nil {
		t.Fatal(err)
	}

	if ws.Name == updateableWs.Name {
		t.Error("Update error at Name")
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
	if err := h.Remove(split[1], split[0], &webshop{}); err != nil {
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
		_, _ = h.Write(newTestStruct1(), "webshop")
	}
}

func BenchmarkGetSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		startInsert := time.Now()
		ID, _ := h.Write(newTestStruct1(), "webshop")
		fmt.Printf("Insert: %vns\tGet: ", time.Since(startInsert).Nanoseconds())
		start := time.Now()
		_ = h.Read(ID, "webshop", webshop{})
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
		_ = h.Read(split[1], split[0], &webshop{})
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
		_ = h.Remove(split[1], split[0], &webshop{})
		fmt.Printf("%vns\n", time.Since(start).Nanoseconds())
	}
}
