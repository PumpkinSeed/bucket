package odatas

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/rs/xid"
	"github.com/volatiletech/null"
	"time"
)

// TestStruct1 is a struct used for testing and represents an order of a webshop
type testStruct1 struct {
	Token                        string
	CreationDate                 string
	ModificationDate             string
	Status                       string
	PaymentMethod                string
	InvoiceNumber                string
	Email                        string
	CardHolderName               string
	CreditCardLast4Digits        string
	BillingAddressName           string
	BillingAddressCompanyName    string
	BillingAddressAddress1       string
	BillingAddressAddress2       string
	BillingAddressCity           string
	BillingAddressCountry        string
	BillingAddressProvince       string
	BillingAddressPostalCode     string
	BillingAddressPhone          string
	Notes                        string
	ShippingAddressName          string
	ShippingAddressCompanyName   string
	ShippingAddressAddress1      string
	ShippingAddressAddress2      string
	ShippingAddressCity          string
	ShippingAddressCountry       string
	ShippingAddressProvince      string
	ShippingAddressPostalCode    string
	ShippingAddressPhone         string
	ShippingAddressSameAsBilling bool
	FinalGrandTotal              int
	ShippingFees                 int
	ShippingMethod               string
	WillBePaidLater              bool
	PaymentTransactionId         string
}

func NewTestStruct1() testStruct1 {
	addr := gofakeit.Address()
	name := gofakeit.Name()
	return testStruct1{
		Token: xid.New().String(),
		CreationDate: time.Now().UTC().String(),
		Status: "processed",
		PaymentMethod: "card",
		InvoiceNumber: fmt.Sprintf("inv-00%d", gofakeit.Number(10000,99999)),
		Email: gofakeit.Email(),
		CardHolderName: name,
		CreditCardLast4Digits: fmt.Sprintf("%d", gofakeit.Number(1000,9999)),
		BillingAddressName: name,
		BillingAddressAddress1: addr.Address,
		BillingAddressCity: addr.City,
		BillingAddressCountry: addr.Country,
		BillingAddressPostalCode: addr.Zip,
		BillingAddressPhone: gofakeit.Phone(),
		Notes: gofakeit.HipsterSentence(10),
		ShippingAddressName: name,
		ShippingAddressAddress1: addr.Address,
		ShippingAddressCity: addr.City,
		ShippingAddressCountry: addr.Country,
		ShippingAddressPostalCode: addr.Zip,
		ShippingAddressPhone: gofakeit.Phone(),
		FinalGrandTotal: 443,
		ShippingFees: 0,
		ShippingMethod: "Free shipping",
		PaymentTransactionId: xid.New().String(),
	}
}

func EmptyString() null.String {
	return null.String{Valid: false}
}

func EmptyBool() null.Bool {
	return null.Bool{Valid: false}
}
