package odatas

import (
	"fmt"
	"github.com/volatiletech/null"
	"testing"
	"time"
)

type TestData struct {

}

func TestSimpleSearchMatch(t *testing.T) {
	placeholderInit()

	for i := 0; i< 10; i++ {
		order := NewTestStruct1()
		_, err := placeholderBucket.Insert("order::"+order.Token, order, 0)
		if err != nil {
			t.Fatal(err)
		}
	}

	fields := []string{
		"Status",
		"PaymentMethod",
		"Email",
		"CardHolderName",
		"BillingAddressName",
		"BillingAddressAddress1",
		"BillingAddressCity",
		"BillingAddressCountry",
		"BillingAddressPostalCode",
		"BillingAddressPhone",
		"Notes",
		"ShippingMethod",
	}
	manager := placeholderBucket.Manager("", "")
	err := manager.CreateIndex("order_idx", fields, true, false)
	if err != nil {
		t.Fatal(err)
	}

	handler := New(&Configuration{})
	searchMatch := "Clay Monahan"
	mes := time.Now()
	err = handler.SimpleSearch("order_idx", &SearchQuery{
		Match: null.StringFrom(searchMatch),
		Field: "CardHolderName",
	})
	fmt.Println(time.Since(mes))
	if err != nil {
		t.Fatal(err)
	}
}