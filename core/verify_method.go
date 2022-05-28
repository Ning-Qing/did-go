package core

import "crypto"

type VerificationMethod struct {
	ID string `json:"id"`

	// Controller 密钥控制器，标识密钥所有者
	Controller   string `json:"controller"`
	Type         string `json:"type"`
	PublicKeyJwk Jwk    `json:"publicKeyJwk"`

	// PublicKeyMultibase 是 [MULTIBASE] 编码公钥的字符串表示形式
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

func NewVerificationMethod(id, controller string, pubKey crypto.PublicKey) *VerificationMethod {
	key := NewPubKeyFactory(pubKey)
	if key == nil {
		return nil
	}
	jwk, t := key.Jwk()
	return &VerificationMethod{
		ID:           id,
		Controller:   controller,
		Type:         t,
		PublicKeyJwk: jwk,
	}
}

func (v *VerificationMethod) Verify(hash []byte, sign []byte) bool {
	pubKey := v.PublicKeyJwk.GeneratePubKey()
	key := NewPubKeyFactory(pubKey)
	if key == nil {
		return false
	}
	return key.Verify(hash, sign)
}
