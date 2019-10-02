package bucket

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

	Opts Opts `json:"bucket_opts"`
}

type Opts struct {
	OperationTimeout      NullTimeout `json:"operation_timeout"`
	BulkOperationTimeout  NullTimeout `json:"bulk_operation_timeout"`
	DurabilityTimeout     NullTimeout `json:"durability_timeout"`
	DurabilityPollTimeout NullTimeout `json:"durability_poll_timeout"`
	ViewTimeout           NullTimeout `json:"view_timeout"`
	N1qlTimeout           NullTimeout `json:"n1ql_timeout"`
	AnalyticsTimeout      NullTimeout `json:"analytics_timeout"`
}

// New creates a new handler from the configuration that handles the operations
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

	h := &Handler{
		http:        client,
		httpAddress: "http://localhost:8091",
		username:    c.Username,
		password:    c.Password,
		state:       s,
	}

	h.prepare()
	return h, nil
}

//GetManager returns a BucketManager for performing management operations on this bucket
func (h *Handler) GetManager(ctx context.Context) *gocb.BucketManager {
	return h.state.bucket.Manager(h.username, h.password)
}

// ValidateState validates the state of the bucket
func (h *Handler) ValidateState() (bool, error) {
	return h.state.validate()
}

func (h *Handler) prepare() {
	h.prepareBucket()
}

func (h *Handler) prepareBucket() {
	if h.state.configuration.Opts.OperationTimeout.valid {
		h.state.bucket.SetOperationTimeout(h.state.configuration.Opts.OperationTimeout.Value)
	}
	if h.state.configuration.Opts.BulkOperationTimeout.valid {
		h.state.bucket.SetBulkOperationTimeout(h.state.configuration.Opts.BulkOperationTimeout.Value)
	}
	if h.state.configuration.Opts.DurabilityTimeout.valid {
		h.state.bucket.SetDurabilityTimeout(h.state.configuration.Opts.DurabilityTimeout.Value)
	}
	if h.state.configuration.Opts.DurabilityPollTimeout.valid {
		h.state.bucket.SetDurabilityPollTimeout(h.state.configuration.Opts.DurabilityPollTimeout.Value)
	}
	if h.state.configuration.Opts.ViewTimeout.valid {
		h.state.bucket.SetViewTimeout(h.state.configuration.Opts.ViewTimeout.Value)
	}
	if h.state.configuration.Opts.N1qlTimeout.valid {
		h.state.bucket.SetN1qlTimeout(h.state.configuration.Opts.N1qlTimeout.Value)
	}
	if h.state.configuration.Opts.AnalyticsTimeout.valid {
		h.state.bucket.SetAnalyticsTimeout(h.state.configuration.Opts.AnalyticsTimeout.Value)
	}
}
