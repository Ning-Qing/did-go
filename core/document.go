package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/json"
	"log"
	"math/big"
)

type Document struct {
	Context []string `json:"@context"`
	ID      string   `json:"id"`

	// AlsoKnownAs 其值为一个或多个 DID
	// 表明此DID也由一个或多个其他DID标识。
	AlsoKnownAs []string `json:"alsoKnownAs"`

	// Controller 其值为一个或多个 DID
	// 这些DID的文档中包含的任何验证方法都可以被接受
	// 因此满足这些验证方法的证明应被视为等同于DID主体提供的证明。
	Controller []string `json:"controller"`

	VerificationMethod []*VerificationMethod `json:"verificationMethod"`
	Service            *Service              `json:"service"`
}

func (d *Document) PutVerifMethod(method *VerificationMethod) {
	d.VerificationMethod = append(d.VerificationMethod, method)
}

func (d *Document) GetVerifMethod(id string) *VerificationMethod {
	for _, v := range d.VerificationMethod {
		if v.ID == id {
			return v
		}
	}
	return nil
}

func (d *Document) Serialization() ([]byte, error) {
	return json.Marshal(d)
}

type VerificationMethod struct {
	ID string `json:"id"`

	// Controller 密钥控制器，标识密钥所有者
	Controller   string      `json:"controller"`
	Type         string      `json:"type"`
	PublicKeyJwk interface{} `json:"publicKeyJwk"`

	// PublicKeyMultibase 是 [MULTIBASE] 编码公钥的字符串表示形式
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type jkw_rsa struct {
	Kty string `json:"kty"`
	E   int    `json:"e"`
	N   string `json:"n"`
}

type jkw_es struct {
	Kty   string `json:"kty"`
	Curve string `json:"curve"`
	X     string `json:"x"`
	Y     string `json:"y"`
}

func CreateVerificationMethod(id, controller string, pubKey crypto.PublicKey) *VerificationMethod {
	switch publicKey := pubKey.(type) {
	case *rsa.PublicKey:
		return &VerificationMethod{
			ID:         id,
			Controller: controller,
			Type:       RSA256,
			PublicKeyJwk: &jkw_rsa{
				Kty: "RSA",
				E:   publicKey.E,
				N:   publicKey.N.String(),
			},
		}
	case *ecdsa.PublicKey:
		return &VerificationMethod{
			ID:         id,
			Controller: controller,
			Type:       ES256,
			PublicKeyJwk: &jkw_es{
				Kty:   "EC",
				Curve: "P256",
				X:     publicKey.X.String(),
				Y:     publicKey.Y.String(),
			},
		}
	default:
		return nil
	}
}

func (v *VerificationMethod) Verify(hash []byte, sign []byte) bool {
	switch v.Type {
	case RSA256:
		jkw := v.PublicKeyJwk.(*jkw_rsa)
		n, _ := new(big.Int).SetString(jkw.N, 10)
		pub := &rsa.PublicKey{
			N: n,
			E: jkw.E,
		}
		err := rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash, sign)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	case ES256:
		jkw := v.PublicKeyJwk.(*jkw_es)
		x, _ := new(big.Int).SetString(jkw.X, 10)
		y, _ := new(big.Int).SetString(jkw.Y, 10)
		curve := elliptic.P256()
		pub := &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}
		return ecdsa.VerifyASN1(pub, hash, sign)
	default:
		log.Println("unsupported type")
		return false
	}
}

type Service struct {
	ID              string   `json:"id"`
	Type            []string `json:"type"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}
