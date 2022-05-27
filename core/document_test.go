package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"testing"
)

func TestCreatDocument(t *testing.T) {
	doc := &Document{
		Context:            []string{"github.com/Ning-Qing"},
		ID:                 "did:example:user",
		VerificationMethod: make([]*VerificationMethod, 0),
	}

	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	verifMethod := CreateVerificationMethod("did:example:user#test", "did:example:issuer", &privKey.PublicKey)
	doc.PutVerifMethod(verifMethod)
	res, _ := doc.Serialization()
	log.Println(string(res))
}

func TestPutVerifMethod(t *testing.T) {
	doc := &Document{
		Context:            []string{"github.com/Ning-Qing"},
		ID:                 "did:example:user",
		VerificationMethod: make([]*VerificationMethod, 0),
	}

	privKey1, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	method1 := CreateVerificationMethod("did:example:user#rsa", "did:example:issuer", &privKey1.PublicKey)
	doc.PutVerifMethod(method1)

	curve := elliptic.P256()
	privKey2, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err)
	}
	method2 := CreateVerificationMethod("did:example:user#es", "did:example:issuer", &privKey2.PublicKey)
	doc.PutVerifMethod(method2)

	res, _ := doc.Serialization()
	log.Println(string(res))
}
