package bucket

const (
	metaFieldName = "_meta"
)

type metaContainer struct {
	Meta *meta `json:"_meta"`
}

type meta struct {
	ReferencedDocuments []referencedDocumentMeta `json:"referenced_documents"`
	Parent              string                   `json:"parent"`
}

type referencedDocumentMeta struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (h *Handler) getMeta(typ, id string) (*meta, error) {
	var c = metaContainer{}
	dk, err := h.state.getDocumentKey(typ, id)
	if err != nil {
		return nil, err
	}
	_, err = h.state.bucket.Get(dk, &c)
	if err != nil {
		return nil, err
	}

	return c.Meta, nil
}

func (m *meta) AddReferencedDocument(key, typ, id string) {
	m.ReferencedDocuments = append(m.ReferencedDocuments, referencedDocumentMeta{
		Key:  key,
		Type: typ,
		ID:   id,
	})
}
