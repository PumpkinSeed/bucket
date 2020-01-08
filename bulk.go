package bucket

import (
	"context"
	"reflect"

	"github.com/couchbase/gocb"
)

// GetBulk accepts a set of hits from a search
// and a container represents the data-structure
// and fill it up with the hits where the container
// should be *[]T type
func (h *Handler) GetBulk(ctx context.Context, hits []gocb.SearchResultHit, container interface{}) error {
	var items []gocb.BulkOp
	rv := reflect.ValueOf(container)
	if rv.Type().Kind() != reflect.Ptr {
		return ErrInvalidBulkContainer
	}

	rvElem := rv.Elem()
	switch rvElem.Kind() {
	case reflect.Slice:
		if rvElem.Len() != len(hits) {
			return ErrInvalidBulkContainer
		}
		for i := 0; i < rvElem.Len(); i++ {
			typs, err := getDocumentTypes(rvElem.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
			items = append(items, &gocb.GetOp{Key: hits[i].Id, Value: rvElem.Index(i).Addr().Interface()})
			identifier := h.state.fetchDocumentIdentifier(hits[i].Id)
			addressableFields := getStructAddressableSubfields(rvElem.Index(i).Addr())
			for _, typ := range typs {
				documentKey := h.state.getDocumentKey(typ, identifier)
				items = append(items, &gocb.GetOp{Key: documentKey, Value: addressableFields[typ]})
			}
		}
	default:
		return ErrInvalidBulkContainer
	}

	return h.state.bucket.Do(items)
}
