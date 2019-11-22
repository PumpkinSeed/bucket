package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	enc.Encode(v)

	if _, err := w.Write(buf.Bytes()); err != nil {
		panic(err)
	}
}