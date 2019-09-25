package odatas

import "gopkg.in/couchbase/gocb.v1"

type Handler struct {
	state *State

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

	s := NewState(b, c.Bucket, c.Separator)
	err = s.Load()
	if err != nil {
		panic(err)
	}

	return Handler{
		bucket: b,
		state:  s,
	}
}
