package bucket

import (
	"context"
	"fmt"
	"reflect"

	"github.com/couchbase/gocb"
)

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
			fmt.Println(typs)
			//identifier := h.state.fetchDocumentIdentifier(hits[i].Id)
			items = append(items, &gocb.GetOp{Key: hits[i].Id, Value: rvElem.Index(i).Addr().Interface()})
			//for _, typ := range typs {
			//	items = append(items, &gocb.GetOp{Key: })
			//}
		}
	default:
		return ErrInvalidBulkContainer
	}

	return h.state.bucket.Do(items)
}
