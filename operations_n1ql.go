package bucket

import "github.com/couchbase/gocb"

func (h *Handler) NewN1qlQuery(statement string) *gocb.N1qlQuery {
	return gocb.NewN1qlQuery(statement)
}

//TODO add custom param
// TODO INNER JOIN our separated documents
// TODO add something like sqlboiler gives us GetName, etc.
// TODO need to reformat query, where FROM will be h.state.bucket
// TODO use parent meta
// reformat query like:
// TODO SELECT name, product.price FROM webshop WHERE id = 12 ---> SELECT name, product.price FROM h.state.bucket JOIN h.state.bucket product ON product.meta().parent = meta().id WHERE id = 12
// TODO what if query starts with SELECT *?
//func (h *Handler) AddParam(params interface{}, p interface{}) interface{} {
//	return append(params, p)
//}

func (h *Handler) ExecuteN1qlQuery(q *gocb.N1qlQuery, params interface{}) (gocb.QueryResults, error) {
	qr, err := h.state.bucket.ExecuteN1qlQuery(q, nil)
	if err != nil {
		return nil, err
	}
	return qr, nil
}
