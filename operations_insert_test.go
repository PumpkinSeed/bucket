package bucket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

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

	wsGet := webshop{}
	if err := th.Get(context.Background(), "webshop", id, &wsGet); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ws, wsGet, "should be equal")

	m, err := th.getMeta("webshop", id)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.ChildDocuments) != 3 {
		t.Errorf("Length of children should be 3, instead of %d", len(m.ChildDocuments))
	}
	if m.Type != "webshop" {
		t.Errorf("Type should be 'webshop', instead of %s", m.Type)
	}
	if m.ParentDocument != nil {
		t.Errorf("Parent should be nil, instead of %v", m.ParentDocument)
	}

	for _, child := range m.ChildDocuments {
		parts := strings.Split(child.Key, "::")
		if len(parts) != 2 {
			t.Error("Invalid key")
		}
		if parts[0] != child.Type {
			t.Errorf("First part of key should be equal with %s, instead of %s", child.Type, parts[0])
		}
		if parts[1] != child.ID {
			t.Errorf("Second part of key should be equal with %s, instead of %s", child.ID, parts[1])
		}
		if id != child.ID {
			t.Errorf("Second part of key should be equal with %s, instead of %s", id, child.ID)
		}
	}

	pGet := product{}
	if err := th.Get(context.Background(), "product", id, &pGet); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, ws.Product, &pGet, "should be equal")

	m, err = th.getMeta("product", id)
	if err != nil {
		t.Fatal(err)
	}

	if len(m.ChildDocuments) != 1 {
		t.Errorf("Length of children should be 1, instead of %d", len(m.ChildDocuments))
	}
	if m.Type != "product" {
		t.Errorf("Type should be 'product', instead of %s", m.Type)
	}
	if m.ParentDocument != nil && m.ParentDocument.Key != "webshop::"+id {
		t.Errorf("Parent key should be %s, instead of %s", "webshop::"+id, m.ParentDocument.Key)
	}

	for _, child := range m.ChildDocuments {
		parts := strings.Split(child.Key, "::")
		if len(parts) != 2 {
			t.Error("Invalid key")
		}
		if parts[0] != child.Type {
			t.Errorf("First part of key should be equal with %s, instead of %s", child.Type, parts[0])
		}
		if parts[1] != child.ID {
			t.Errorf("Second part of key should be equal with %s, instead of %s", child.ID, parts[1])
		}
		if id != child.ID {
			t.Errorf("Second part of key should be equal with %s, instead of %s", id, child.ID)
		}
	}
}
