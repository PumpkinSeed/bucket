package odatas

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/rs/xid"
)

const (
	bucketName = "company"
)

var th *Handler
var seeded bool

func init() {
	seed()
}

func seed() {
	if th == nil {
		th = defaultHandler()
	}

	th.SetDocumentType(context.Background(), "order", "order")
	th.SetDocumentType(context.Background(), "webshop", "webshop")
	th.SetDocumentType(context.Background(), "product", "product")

	var test = os.Getenv("PKG_TEST")
	if test == "testing" && !seeded {
		log.Print("TestEnv")
		gofakeit.Seed(time.Now().UnixNano())

		start := time.Now()
		if err := th.GetManager(context.Background()).Flush(); err != nil {
			fmt.Printf("Turn on flush in bucket: %+v\n", err)
		}
		fmt.Printf("Bucket flushed: %v\n", time.Since(start))

		for j := 0; j < 1000; j++ {
			instance := generate()
			th.Insert(context.Background(), "webshop", "", instance)
			// _, _ = th.state.bucket.Insert(instance.Token, instance, 0)
		}
		fmt.Printf("Connection setup, data seeded %v\n", time.Since(start))
		seeded = true
	}
}

// webshop is a struct used for testing and represents an order of a webshop
type webshop struct {
	Token                        string   `json:"token"`
	CreationDate                 string   `json:"creation_date"`
	ModificationDate             string   `json:"modification_date"`
	Status                       string   `json:"status"`
	PaymentMethod                string   `json:"payment_method"`
	InvoiceNumber                string   `json:"invoice_number"`
	Email                        string   `json:"email" indexable:"true"`
	CardHolderName               string   `json:"card_holder_name"`
	CreditCardLast4Digits        string   `json:"credit_card_last_4_digits"`
	BillingAddressName           string   `json:"billing_address_name" indexable:"true"`
	BillingAddressCompanyName    string   `json:"billing_address_company_name" indexable:"true"`
	BillingAddressAddress1       string   `json:"billing_address_address_1"`
	BillingAddressAddress2       string   `json:"billing_address_address_2"`
	BillingAddressCity           string   `json:"billing_address_city"`
	BillingAddressCountry        string   `json:"billing_address_country"`
	BillingAddressProvince       string   `json:"billing_address_province"`
	BillingAddressPostalCode     string   `json:"billing_address_postal_code"`
	BillingAddressPhone          string   `json:"billing_address_phone"`
	Notes                        string   `json:"notes"`
	ShippingAddressName          string   `json:"shipping_address_name"`
	ShippingAddressCompanyName   string   `json:"shipping_address_company_name"`
	ShippingAddressAddress1      string   `json:"shipping_address_address_1"`
	ShippingAddressAddress2      string   `json:"shipping_address_address_2"`
	ShippingAddressCity          string   `json:"shipping_address_city"`
	ShippingAddressCountry       string   `json:"shipping_address_country"`
	ShippingAddressProvince      string   `json:"shipping_address_province"`
	ShippingAddressPostalCode    string   `json:"shipping_address_postal_code"`
	ShippingAddressPhone         string   `json:"shipping_address_phone"`
	ShippingAddressSameAsBilling bool     `json:"shipping_address_same_as_billing"`
	FinalGrandTotal              int      `json:"final_grand_total"`
	ShippingFees                 int      `json:"shipping_fees"`
	ShippingMethod               string   `json:"shipping_method"`
	WillBePaidLater              bool     `json:"will_be_paid_later"`
	PaymentTransactionId         string   `json:"payment_transaction_id"`
	Product                      *product `json:"product"`
	Store                        *store   `json:"store,omitempty"`
}

type product struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	StoreID     string `json:"store_id"`
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
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func generate() webshop {
	addr := gofakeit.Address()
	name := gofakeit.Name()
	beer := gofakeit.BeerName()
	return webshop{
		Token:                     xid.New().String(),
		CreationDate:              time.Now().UTC().String(),
		Status:                    "processed",
		PaymentMethod:             "card",
		InvoiceNumber:             fmt.Sprintf("inv-00%d", gofakeit.Number(10000, 99999)),
		Email:                     gofakeit.Email(),
		CardHolderName:            name,
		CreditCardLast4Digits:     fmt.Sprintf("%d", gofakeit.Number(1000, 9999)),
		BillingAddressName:        name,
		BillingAddressAddress1:    addr.Address,
		BillingAddressCity:        addr.City,
		BillingAddressCountry:     addr.Country,
		BillingAddressPostalCode:  addr.Zip,
		BillingAddressPhone:       gofakeit.Phone(),
		Notes:                     gofakeit.HipsterSentence(10),
		ShippingAddressName:       name,
		ShippingAddressAddress1:   addr.Address,
		ShippingAddressCity:       addr.City,
		ShippingAddressCountry:    addr.Country,
		ShippingAddressPostalCode: addr.Zip,
		ShippingAddressPhone:      gofakeit.Phone(),
		FinalGrandTotal:           443,
		ShippingFees:              0,
		ShippingMethod:            "Free shipping",
		PaymentTransactionId:      xid.New().String(),
		Product: &product{
			ID:          xid.New().String(),
			UserID:      xid.New().String(),
			StoreID:     xid.New().String(),
			Name:        beer,
			Description: gofakeit.HipsterSentence(10),
			Slug:        strings.ToLower(beer),
			Price:       1233,
			SalePrice:   1400,
			CurrencyID:  2,
			OnSale:      123,
			Status:      "active",
		},
		Store: &store{
			ID:          xid.New().String(),
			UserID:      xid.New().String(),
			Name:        "productshop",
			Description: "Product shop",
		},
	}

}
