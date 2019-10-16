package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestHandler_GetSubDocuments(t *testing.T) {
	var ws = generate()

	resultset := th.getSubDocuments("webshop", xid.New().String(), ws, nil)
	for k, v := range resultset {
		shit, _ := json.Marshal(v)
		fmt.Printf("k: %s, v: %s\n", k, shit)
	}
}

func TestHandler_Insert(t *testing.T) {
	ws := generate()
	_, id, err := th.EInsert(context.Background(), "webshop", "", ws, 0)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(id)
}
