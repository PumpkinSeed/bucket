package odatas

import (
	"testing"

	"github.com/couchbase/gocb"
)

const (
	bucketName = "company"
)

var i *indexer

type basicUser struct {
	ID          string `json:"id" indexable:"true"`
	Name        string `json:"name" indexable:"true"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Password    string `json:"-"`
}

func init() {
	cluster, _ := gocb.Connect("couchbase://localhost")
	_ = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Administrator",
		Password: "password",
	})
	bucket, _ := cluster.OpenBucket(bucketName, "")
	i = NewIndexer(bucket, "Administrator", "password")
}

func TestIndex(t *testing.T) {
	instance := basicUser{}

	if err := i.Index(instance); err != nil {
		t.Fatal(err)
	}

	indexes, err := i.bucketManager.GetIndexes()
	if err != nil {
		t.Fatal(err)
	}

	if len(indexes) < 2 {
		t.Error("Missing indexes")
	}

	for _, ind := range indexes {
		t.Logf("%+v", ind.Name)
	}

}
