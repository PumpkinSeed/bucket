package bucket

import (
	"strings"
	"testing"
)

func TestGetMeta(t *testing.T) {
	_, id, err := testInsert()
	if err != nil || id == "" {
		t.Fatal(err)
	}

	m, err := th.getMeta("webshop", id)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(m.ReferencedDocuments[0].Key, "origin::") {
		t.Errorf("Referenced first elem should contain 'origin::', instead of %s", m.ReferencedDocuments[0])
	}
}
