package bucket

import (
	"reflect"
	"testing"

	"github.com/rs/xid"
)

func TestGetDocumentTypesWithPointer(t *testing.T) {
	typs, err := getDocumentTypes(&webshop{})
	if err != nil {
		t.Fatal(err)
	}
	var expected = []string{"product", "store"}
	if !reflect.DeepEqual(typs, expected) {
		t.Errorf("Types should be %v, instead of %v", expected, typs)
	}
}

func TestGetDocumentTypes(t *testing.T) {
	_, err := getDocumentTypes(webshop{})
	if err != ErrInvalidGetDocumentTypesParam {
		t.Errorf("Error should be %v, instead of %v", ErrInvalidGetDocumentTypesParam, err)
	}
}

func TestGetDocumentTypesNonPointerExpectError(t *testing.T) {
	v := "webshop"
	if _, err := getDocumentTypes(&v); err == nil {
		t.Errorf("Error should be value argument must be a struct instead of nil")
	}
}

func TestGetStructAddressableSubfields(t *testing.T) {
	var ws = &webshop{}
	var s = store{
		ID:          xid.New().String(),
		UserID:      xid.New().String(),
		Name:        "test",
		Description: "test description",
	}

	result := getStructAddressableSubfields(reflect.ValueOf(ws))
	result["store"] = &s
	if v, ok := result["store"].(*store); ok {
		if v.Name != "test" {
			t.Errorf("Name should be 'test', instead of %s", v.Name)
		}
	}
}
