package bucket

import (
	"context"

	"github.com/couchbase/gocb"
)

type writerF func(string, string, interface{}, uint32) (gocb.Cas, error)

//type readerF func(string, string, interface{}, uint32) (gocb.Cas, error)

// Cas is the container of Cas operation of all documents
type Cas map[string]gocb.Cas

// Remove removes a document from the bucket
func (h *Handler) Remove(ctx context.Context, typ, id string, ptr interface{}) error {
	typs, e := getDocumentTypes(ptr)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		documentID := h.state.getDocumentKey(typ, id)
		if _, err := h.state.bucket.Remove(documentID, 0); err != nil {
			return err
		}
	}
	return nil
}

// Touch touches documents, specifying a new expiry time for it
// The Cas value must be 0
func (h *Handler) Touch(ctx context.Context, typ, id string, ptr interface{}, ttl uint32) error {
	typs, e := getDocumentTypes(ptr)
	if e != nil {
		return e
	}

	for _, typ := range typs {
		documentID := h.state.getDocumentKey(typ, id)
		if _, err := h.state.bucket.Touch(documentID, 0, ttl); err != nil {
			return err
		}
	}
	return nil
}

// Ping will ping a list of services and verify they are active and responding in an acceptable period of time
func (h *Handler) Ping(ctx context.Context, services []gocb.ServiceType) (*gocb.PingReport, error) {
	report, err := h.state.bucket.Ping(services)
	if err != nil {
		return nil, err
	}
	return report, nil
}
