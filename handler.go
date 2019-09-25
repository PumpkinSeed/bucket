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

func New(c *Configuration) (*Handler, error) {
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
		return nil, err
	}
	if err := cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	}); err != nil {
		return nil, err
	}
	bucket, err := cluster.OpenBucket(c.BucketName, "")
	if err != nil {
		return nil, err
	}

	return &Handler{
		http:        client,
		httpAddress: "http://localhost:8091",
		bucket:      bucket,
		username:    c.Username,
		password:    c.Password,
	}, nil
}

func (h *Handler) GetManager() *gocb.BucketManager {
	return h.bucket.Manager(h.username, h.password)
}
