package core

import "sync"

type Provider interface {
	GetDocument(id string) *Document
	PutDocument(doc *Document)
	UpdataDocument(doc *Document)
	DeleteDocument(id string)
}

var _ Provider = (*MemoryProvider)(nil)

type MemoryProvider struct {
	docs sync.Map
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{}
}

func (p *MemoryProvider) GetDocument(id string) *Document {
	v, ok := p.docs.Load(id)
	if !ok {
		return nil
	}
	doc, ok := v.(Document)
	if !ok {
		return nil
	}
	return &doc
}

func (p *MemoryProvider) PutDocument(doc *Document) {
	p.docs.Store(doc.ID, *doc)
}

func (p *MemoryProvider) UpdataDocument(doc *Document) {
	p.docs.Store(doc.ID, *doc)
}

func (p *MemoryProvider) DeleteDocument(id string) {
	p.docs.Delete(id)
}
