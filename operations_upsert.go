package bucket

import (
	"context"

	"github.com/couchbase/gocb"
	"github.com/rs/xid"
)

// Upsert inserts or replaces a document in the bucket
func (h *Handler) Upsert(ctx context.Context, typ, id string, q interface{}, ttl uint32) (Cas, string, error) {
	if id == "" {
		id = xid.New().String()
	}

	kv := h.getSubDocuments(typ, id, q, nil)

	var ops []gocb.BulkOp
	for k, v := range kv {
		key := h.state.getDocumentKey(k, id)
		ops = append(ops, &gocb.UpsertOp{Key: key, Value: v, Expiry: ttl})
	}

	err := h.state.bucket.Do(ops)
	return nil, id, err
}
