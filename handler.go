package odatas

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/couchbase/gocb"
)

type Handler struct {
	state state

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
	Separator        string `json:"separator"`
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
	bucket, _ := cluster.OpenBucket(c.BucketName, c.BucketPassword)

	s := newState(bucket, c.Separator)
	_ = s.load()

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

func (h *Handler) Insert(ctx *context.Context, bucketName, key string, value interface{}) error {
	prefix, err := h.state.getType(bucketName)
	if err != nil {
		if err2 := h.state.newType(bucketName, bucketName); err2 != nil {
			return err2
		}

		prefix = bucketName + h.state.separator
	}
	if _, err = h.bucket.Upsert(prefix+key, value, 0); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Get(ctx *context.Context, bucketName, key string) (interface{}, error) {
	var res interface{}
	prefix, err := h.state.getType(bucketName)
	if err != nil {
		return res, err
	}

	_, err = h.bucket.Get(prefix+key, res)
	if err != nil {
		return res, err
	}

	return res, nil
}
