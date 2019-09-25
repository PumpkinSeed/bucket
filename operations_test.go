package odatas

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	placeholderInit()
	ws := webshop{
		ID:       23,
		RoleID:   25,
		Name:     "Test",
		Email:    "test@test.com",
		Password: "asd",
		Product:  product{
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
		Store:    store{
			ID:          55,
			UserID:      44,
			Name:        "productshop",
			Description: "Product shop",
		},
	}
	err := Insert(ws, "webshop")
	if err != nil {
		t.Fatal(err)
	}
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
	placeholderInit()
	ws := webshop{
		ID:       0,
		RoleID:   0,
		Name:     "",
		Email:    "",
		Password: "",
		Product:  product{
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
		Store:    store{
			ID:          0,
			UserID:      0,
			Name:        "",
			Description: "",
		},
	}
	err := read("bm5lf7cudj7ugmucbro0","webshop",&ws )
	if err != nil {
		t.Fail()
	}
	fmt.Printf("%+v\n", ws)
}
