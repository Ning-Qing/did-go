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

func TestIssuerClaimRSA256(t *testing.T) {
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
	res, _ := claim.Serialization()
	log.Println(string(res))
}

func TestIssuerClaimES256(t *testing.T) {
	curve := elliptic.P256()
	privKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err)
	}
	claim := &Claim{
		Context: []string{"github.com/Ning-Qing"},
		ID:      "did:example:user",
		Doc: &Doc{
			Issuer:         "did:example:issuer",
			IssuanceDate:   time.Now().String(),
			ExpirationDate: time.Now().AddDate(1, 0, 0).String(),
			Revocation:     "github.com/Ning-Qing",
		},
	}
	claim.Sign("did:example:issuer", "did:example:user#did", privKey)
	res, _ := claim.Serialization()
	log.Println(string(res))
}
