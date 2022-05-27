package core

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"testing"
	"time"
)

func TestVCRSA256(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	claim := &Claim{
		Context: []string{"github.com/Ning-Qing"},
		ID:      "did:example:user",
		Doc: &Doc{
			Issuer:            "did:example:issuer",
			IssuanceDate:      time.Now().String(),
			ExpirationDate:    time.Now().AddDate(1, 0, 0).String(),
			Revocation:        "github.com/Ning-Qing",
			CredentialSubject: nil,
		},
	}
	claim.Sign("did:example:issuer", "did:example:user#test", privKey)

	doc := &Document{
		Context:            []string{"github.com/Ning-Qing"},
		ID:                 "did:example:user",
		VerificationMethod: make([]*VerificationMethod, 0),
	}

	method := CreateVerificationMethod("did:example:user#test", "did:example:issuer", &privKey.PublicKey)
	doc.PutVerifMethod(method)

	log.Println(ValidateClaim(claim, doc))
}
