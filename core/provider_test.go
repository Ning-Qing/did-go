package core

import (
	"log"
	"testing"
)

func TestMemoryProvider(t *testing.T) {

	provider := NewMemoryProvider()

	doc := &Document{
		Context:            []string{"github.com/Ning-Qing"},
		ID:                 "did:example:user",
		VerificationMethod: make([]*VerificationMethod, 0),
	}

	provider.PutDocument(doc)

	doc.Service = &Service{
		ID:              "did:example:user1",
		Type:            []string{"test"},
		ServiceEndpoint: "127.0.0.1",
	}

	provider.UpdataDocument(doc)

	doc_updata := provider.GetDocument("did:example:user")
	res, _ := doc_updata.Serialization()
	log.Println(string(res))

	provider.DeleteDocument("did:example:user")

	log.Println(provider.GetDocument("did:example:user"))
}
