package bucket

const (
	metaFieldName = "_meta"
)

type metaContainer struct {
	Meta *meta `json:"_meta"`
}

type meta struct {
	ChildDocuments []documentMeta `json:"_children"`
	ParentDocument *documentMeta  `json:"_parent"`
	Type           string         `json:"_type"`
}

type documentMeta struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	ID   string `json:"id"`
}

func (h *Handler) getMeta(typ, id string) (*meta, error) {
	var c = metaContainer{}
	dk := h.state.getDocumentKey(typ, id)

	_, err := h.state.bucket.Get(dk, &c)
	if err != nil {
		return nil, err
	}

	return c.Meta, nil
}

func (m *meta) AddChildDocument(key, typ, id string) {
	m.ChildDocuments = append(m.ChildDocuments, documentMeta{
		Key:  key,
		Type: typ,
		ID:   id,
	})
}
