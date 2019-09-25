package odatas

import (
	"gopkg.in/couchbase/gocb.v1"
	"net"
	"net/http"
	"time"

)

type Handler struct {
	state string

	address     string
	httpAddress string

	http *http.Client

	bucket *gocb.Bucket

	username string // temp field
	password string // temp field
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

	cluster, err := gocb.Connect(c.ConnectionString)
	if err != nil {
		panic(err)
	}
	if err := cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	}); err != nil {
		panic(err)
	}
	bucket, err := cluster.OpenBucket(c.BucketName, "")
	if err != nil {
		panic(err)
	}

	return Handler{
		http:        client,
		httpAddress: "http://localhost:8091",
		bucket:      bucket,
		username:    c.Username,
		password:    c.Password,
	}
}

func (h *Handler) GetManager() *gocb.BucketManager {
	return h.bucket.Manager(h.username, h.password)
}
