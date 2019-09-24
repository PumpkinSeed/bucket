package odatas

import (
	"net"
	"net/http"
	"time"
)

type Handler struct {
	state string

	address     string
	httpAddress string

	http *http.Client
}

type Configuration struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	BucketUsername string `json:"bucket_username"`
	BucketPassword string `json:"bucket_password"`
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

	return Handler{
		http: client,
		httpAddress: "http://localhost:8091",
	}
}


