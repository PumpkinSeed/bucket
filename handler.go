package odatas

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/couchbase/gocb"
)

type Handler struct {
	state *state

	//address     string // not used
	httpAddress string

	http *http.Client

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

	s, err := newState(c)
	if err != nil {
		return nil, err
	}
	_ = s.load()

	return &Handler{
		http:        client,
		httpAddress: "http://localhost:8091",
		username:    c.Username,
		password:    c.Password,
		state:       s,
	}, nil
}

func (h *Handler) GetManager(ctx context.Context) *gocb.BucketManager {
	return h.state.bucket.Manager(h.username, h.password)
}

func (h *Handler) ValidateState() (bool, error) {
	return h.state.validate()
}
