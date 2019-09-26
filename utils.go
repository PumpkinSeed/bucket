package odatas

import (
	"encoding/base64"
	"log"
	"net/http"
)

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func setupBasicAuth(req *http.Request) {
	req.Header.Add("Authorization", "Basic "+basicAuth("Administrator", "password"))
}

func defaultHandler() *Handler {
	h, err := New(&Configuration{
		Username:       "Administrator",
		Password:       "password",
		BucketName:     bucketName,
		BucketPassword: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return h
}
