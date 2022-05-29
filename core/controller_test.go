package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"log"
	"testing"
	"time"
)

func TestNewDID(t *testing.T) {
	provider := NewMemoryProvider()
	controller := NewController(provider)
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	pubKey := privKey.PublicKey

	did := controller.NewDID(NewJwkRSA(&pubKey))

	doc := controller.provider.GetDocument(did)
	res, _ := doc.Serialization()
	log.Println(string(res))
}

func TestGenerate(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	log.Println(generate(&privKey.PublicKey))
}

func TestVCRSA256(t *testing.T) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Println(err)
	}
	claim := &Claim{
		Context: []string{"github.com/Ning-Qing"},
		Doc: &Doc{
			Issuer:            "did:example:issuer",
			IssuanceDate:      time.Now().String(),
			ExpirationDate:    time.Now().AddDate(1, 0, 0).String(),
			Revocation:        "github.com/Ning-Qing",
			CredentialSubject: nil,
		},
	}
	claim.Sign("did:example:user#test", privKey)

	method := NewVerificationMethod("did:example:user#test", "did:example:issuer", &privKey.PublicKey)

	log.Println(validateClaim(claim, method))
}

func TestVCES256(t *testing.T) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err)
	}
	claim := &Claim{
		Context: []string{"github.com/Ning-Qing"},
		Doc: &Doc{
			Issuer:            "did:example:issuer",
			IssuanceDate:      time.Now().String(),
			ExpirationDate:    time.Now().AddDate(1, 0, 0).String(),
			Revocation:        "github.com/Ning-Qing",
			CredentialSubject: nil,
		},
	}
	claim.Sign("did:example:user#test", privKey)

	method := NewVerificationMethod("did:example:user#test", "did:example:issuer", &privKey.PublicKey)

	log.Println(validateClaim(claim, method))
}
