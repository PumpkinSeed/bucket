package odatas

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/rs/xid"
)

// TestStruct1 is a struct used for testing and represents an order of a webshop
type testStruct1 struct {
	Token                        string `json:"token"`
	CreationDate                 string `json:"creation_date"`
	ModificationDate             string `json:"modification_date"`
	Status                       string `json:"status"`
	PaymentMethod                string `json:"payment_method"`
	InvoiceNumber                string `json:"invoice_number"`
	Email                        string `json:"email" indexable:"true"`
	CardHolderName               string `json:"card_holder_name"`
	CreditCardLast4Digits        string `json:"credit_card_last_4_digits"`
	BillingAddressName           string `json:"billing_address_name" indexable:"true"`
	BillingAddressCompanyName    string `json:"billing_address_company_name" indexable:"true"`
	BillingAddressAddress1       string `json:"billing_address_address_1"`
	BillingAddressAddress2       string `json:"billing_address_address_2"`
	BillingAddressCity           string `json:"billing_address_city"`
	BillingAddressCountry        string `json:"billing_address_country"`
	BillingAddressProvince       string `json:"billing_address_province"`
	BillingAddressPostalCode     string `json:"billing_address_postal_code"`
	BillingAddressPhone          string `json:"billing_address_phone"`
	Notes                        string `json:"notes"`
	ShippingAddressName          string `json:"shipping_address_name"`
	ShippingAddressCompanyName   string `json:"shipping_address_company_name"`
	ShippingAddressAddress1      string `json:"shipping_address_address_1"`
	ShippingAddressAddress2      string `json:"shipping_address_address_2"`
	ShippingAddressCity          string `json:"shipping_address_city"`
	ShippingAddressCountry       string `json:"shipping_address_country"`
	ShippingAddressProvince      string `json:"shipping_address_province"`
	ShippingAddressPostalCode    string `json:"shipping_address_postal_code"`
	ShippingAddressPhone         string `json:"shipping_address_phone"`
	ShippingAddressSameAsBilling bool   `json:"shipping_address_same_as_billing"`
	FinalGrandTotal              int    `json:"final_grand_total"`
	ShippingFees                 int    `json:"shipping_fees"`
	ShippingMethod               string `json:"shipping_method"`
	WillBePaidLater              bool   `json:"will_be_paid_later"`
	PaymentTransactionId         string `json:"payment_transaction_id"`
}

type testStructEmbedded struct {
	BasicData testStruct1 `json:"basic_data"`
	Tags      []string    `json:"tags"`
	CreatedAt time.Time   `json:"created_at" indexable:"true"`
	DeletedAt time.Time   `json:"deleted_at" indexable:"true"`
}

func newTestStruct1() testStruct1 {
	addr := gofakeit.Address()
	name := gofakeit.Name()
	return testStruct1{
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
	}
}

func emptyString() string {
	return ""
}

func emptyBool() bool {
	return false
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func setupBasicAuth(req *http.Request) {
	req.Header.Add("Authorization", "Basic "+basicAuth("Administrator", "password"))
}
