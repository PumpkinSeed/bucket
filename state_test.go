package bucket

import "testing"

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
