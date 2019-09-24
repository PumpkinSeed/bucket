package odatas

import (
	"net"
	"net/http"
	"time"

	"github.com/couchbase/gocb"
)

type Handler struct {
	state string

	address     string
	httpAddress string

	http *http.Client

	bucket        *gocb.Bucket
	bucketManager *gocb.BucketManager
}

type Configuration struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	BucketName       string `json:"bucket_name"`
	BucketPassword   string `json:"bucket_password"`
	ConnectionString string `json:"connection_string"`
}

func New(c *Configuration) Handler {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: time.Second * 10,
	}

	cluster, _ := gocb.Connect(c.ConnectionString)
	_ = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	})
	bucket, _ := cluster.OpenBucket(c.BucketName, "")

	return Handler{
		http:          client,
		httpAddress:   "http://localhost:8091",
		bucket:        bucket,
		bucketManager: bucket.Manager(c.Username, c.Password),
	}
}
