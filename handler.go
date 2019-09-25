package odatas

import (
	"context"
	
	"gopkg.in/couchbase/gocb.v1"
)

type Handler struct {
	state *state

	bucket *gocb.Bucket
}

type Configuration struct {
	Connection string `json:"connection"`

	Username       string `json:"username"`
	Password       string `json:"password"`
	Bucket         string `json:"bucket"`
	BucketUsername string `json:"bucket_username"`
	BucketPassword string `json:"bucket_password"`
	Separator      string `json:"separator"`
}

func New(c *Configuration) Handler {
	cdb, err := gocb.Connect(c.Connection)
	if err != nil {
		panic(err)
	}

	err = cdb.Authenticate(gocb.PasswordAuthenticator{
		Username: c.Username,
		Password: c.Password,
	})
	if err != nil {
		panic(err)
	}

	b, err := cdb.OpenBucket(c.Bucket, c.BucketPassword)
	if err != nil {
		panic(err)
	}

	s := newState(b, c.Separator)
	err = s.load()
	if err != nil {
		panic(err)
	}

	return Handler{
		bucket: b,
		state:  s,
	}
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

