package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
)

const (
	RSA256 = "RSA256_PKSC1-V1.5_2048_SHA-256"
	ES256  = "ECDSA_P256_SHA-256"
)

type PubKeyFactory interface {
	Jwk() (Jwk, string)
	Verify(hash []byte, sign []byte) bool
	Marshal() ([]byte, error)
}

func NewPubKeyFactory(pubKey crypto.PublicKey) PubKeyFactory {
	switch publicKey := pubKey.(type) {
	case *rsa.PublicKey:
		return &RSAKey{
			pubKey: publicKey,
		}
	case *ecdsa.PublicKey:
		return &ESKey{
			pubKey: publicKey,
		}
	default:
		return nil
	}
}

type PrivKeyFactory interface {
	Sign([]byte) ([]byte, error)
}

func NewPrivKeyFactory(privKey crypto.PrivateKey) PrivKeyFactory {
	switch privKey := privKey.(type) {
	case *rsa.PrivateKey:
		return &RSAKey{
			privKey: privKey,
		}
	case *ecdsa.PrivateKey:
		return &ESKey{
			privKey: privKey,
		}
	default:
		return nil
	}
}

type RSAKey struct {
	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

func (k *RSAKey) Jwk() (Jwk, string) {
	return NewJwkRSA(k.pubKey), RSA256
}

func (k *RSAKey) Verify(hash []byte, sign []byte) bool {
	err := rsa.VerifyPKCS1v15(k.pubKey, crypto.SHA256, hash, sign)
	return err == nil
}

func (k *RSAKey) Marshal() ([]byte, error) {
	return json.Marshal(k.pubKey)
}

func (k *RSAKey) Sign(digest []byte) ([]byte, error) {
	sign, err := k.privKey.Sign(rand.Reader, digest, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

type ESKey struct {
	pubKey  *ecdsa.PublicKey
	privKey *ecdsa.PrivateKey
}

func (k *ESKey) Jwk() (Jwk, string) {
	return NewJwkES(k.pubKey), ES256
}

func (k *ESKey) Verify(hash []byte, sign []byte) bool {
	return ecdsa.VerifyASN1(k.pubKey, hash, sign)
}

func (k *ESKey) Marshal() ([]byte, error) {
	return json.Marshal(k.pubKey)
}

func (k *ESKey) Sign(digest []byte) ([]byte, error) {
	sign, err := k.privKey.Sign(rand.Reader, digest, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	return sign, nil
}
