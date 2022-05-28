package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"math/big"
)

type Jwk interface {
	GeneratePubKey() crypto.PublicKey
}

var _ Jwk = (*JwkRSA)(nil)

type JwkRSA struct {
	Kty string `json:"kty"`
	E   int    `json:"e"`
	N   string `json:"n"`
}

func NewJwkRSA(pubKey *rsa.PublicKey) *JwkRSA {
	return &JwkRSA{
		Kty: "RSA",
		E:   pubKey.E,
		N:   pubKey.N.String(),
	}
}

func (jwk *JwkRSA) GeneratePubKey() crypto.PublicKey {
	n, _ := new(big.Int).SetString(jwk.N, 10)
	return &rsa.PublicKey{
		N: n,
		E: jwk.E,
	}
}

var _ Jwk = (*JwkES)(nil)

type JwkES struct {
	Kty   string `json:"kty"`
	Curve string `json:"curve"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

func NewJwkES(pubKey *ecdsa.PublicKey) *JwkES {
	return &JwkES{
		Kty:   "EC",
		Curve: "P256",
		X:     pubKey.X.String(),
		Y:     pubKey.Y.String(),
	}
}

func (jwk *JwkES) GeneratePubKey() crypto.PublicKey {
	x, _ := new(big.Int).SetString(jwk.X, 10)
	y, _ := new(big.Int).SetString(jwk.Y, 10)
	curve := elliptic.P256()
	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}
}
