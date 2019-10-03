package bucket

import (
	"github.com/rs/xid"
	"testing"
)

func TestUpdateState(t *testing.T) {
	s, err := newState(&Configuration{
		Username:         "Administrator",
		Password:         "password",
		BucketName:       "company",
		BucketPassword:   "",
		ConnectionString: "couchbase://localhost",
		Separator:        "::",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = s.setType("cache", "cache")
	if err != nil {
		t.Fatal(err)
	}

	var s2 state
	_, err = s.bucket.Get(stateDocumentKey, &s2)
	if err != nil {
		t.Fatal(err)
	}

	if s2.DocumentTypes["cache"] != "cache::" {
		t.Errorf("Document key should be 'cache::', instead of %s", s2.DocumentTypes["cache"])
	}
}

func TestValidate(t *testing.T) {
	_, err := th.ValidateState()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	if err := th.state.load(); err != nil {
		t.Error(err)
	}
}

func TestInspect(t *testing.T) {
	_ = th.state.setType("tsinspect", "tsinspect")
	if b := th.state.inspect("tsinspect"); !b {
		t.Error("type not found")
	}
	_ = th.state.deleteType("webshop")

}

func TestSetType(t *testing.T) {
	if err := th.state.setType("webshop", "webshop"); err != nil {
		t.Error(err)
	}
}

func TestDeleteType(t *testing.T) {
	err := th.state.deleteType("webshop")
	if err != nil {
		t.Error(err)
	}
	_ = th.state.setType("webshop", "webshop")
}

func TestFetchDocIdentifierEmptyDocumentKey(t *testing.T) {
	if s := th.state.fetchDocumentIdentifier(""); s != "" {
		t.Errorf("error should be %s instead of %s", "", s)
	}
}

func TestFetchDocIdentifier(t *testing.T) {
	id := xid.New().String()
	if s := th.state.fetchDocumentIdentifier("ws::" + id); s != id {
		t.Errorf("error should be %s instead of %s", id, s)
	}
}
