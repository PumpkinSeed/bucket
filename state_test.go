package bucket

import (
	"testing"

	"github.com/rs/xid"
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

	s.setType("cache", "cache")

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

func TestValidateExpectError(t *testing.T) {
	_ = th.state.updateState()
	delete(th.state.DocumentTypes, "webshop")
	if _, err := th.state.validate(); err == nil {
		t.Error("Err should be not nil")
	}
	_ = th.state.updateState()

}

func TestLoad(t *testing.T) {
	_ = th.state.updateState()
	if err := th.state.load(); err != nil {
		t.Error(err)
	}
}

func TestInspect(t *testing.T) {
	th.state.setType("tsinspect", "tsinspect")
	if b := th.state.inspect("tsinspect"); !b {
		t.Error("type not found")
	}
	_ = th.state.deleteType("webshop")

}

func TestSetType(t *testing.T) {
	th.state.setType("webshop", "webshop")
}

func TestDeleteType(t *testing.T) {
	err := th.state.deleteType("webshop")
	if err != nil {
		t.Error(err)
	}
	th.state.setType("webshop", "webshop")
}

func TestDeleteTypeErrDocumentTypeDoesntExists(t *testing.T) {
	err := th.state.deleteType("webshop")
	if err != nil {
		t.Error(err)
	}
	err = th.state.deleteType("webshop")
	if err != ErrDocumentTypeDoesntExists {
		t.Error(err)
	}
	th.state.setType("webshop", "webshop")
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

func TestCreateEmptyIndexExpectError(t *testing.T) {
	if err := createFullTextSearchIndex("", true, "webshop"); err != ErrEmptyIndex {
		t.Errorf("error should be %s instead of %s", ErrEmptyIndex, err)
	}
}
