package core

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Doc struct {
	Issuer            string                 `json:"issuer"`
	IssuanceDate      string                 `json:"issuanceDate"`
	ExpirationDate    string                 `json:"expirationDate"`
	Revocation        string                 `json:"revocation"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
}

type Claim struct {
	Context []string `json:"@context"`
	ID      string   `json:"id"`
	*Doc
	Proof *Proof `json:"proof"`
}

type Proof struct {
	Creator            string `json:"creator"`
	Type               string `json:"type"`
	SignatureValue     string `json:"signatureValue"`
	VerificationMethod string `json:"verificationMethod"`
}

func (c *Claim) Serialization() ([]byte, error) {
	return json.Marshal(c)
}

func (c *Claim) Hash() ([]byte, error) {
	data, err := json.Marshal(c.Doc)
	if err != nil {
		return nil, err
	}
	hash := sha256.New()
	hash.Write(data)

	return hash.Sum(nil), nil
}

func (c *Claim) Sign(signer, method string, privkey crypto.PrivateKey) error {
	switch privateKey := privkey.(type) {
	case *rsa.PrivateKey:
		digest, err := c.Hash()
		if err != nil {
			return err
		}
		sign, err := privateKey.Sign(rand.Reader, digest, crypto.SHA256)
		if err != nil {
			return err
		}
		signatureValue := base64.StdEncoding.EncodeToString(sign)
		c.Proof = &Proof{
			Creator:            signer,
			Type:               RSA256,
			SignatureValue:     signatureValue,
			VerificationMethod: method,
		}
		return nil
	case *ecdsa.PrivateKey:
		digest, err := c.Hash()
		if err != nil {
			return err
		}
		sign, err := privateKey.Sign(rand.Reader, digest, crypto.SHA256)
		if err != nil {
			return err
		}
		signatureValue := base64.StdEncoding.EncodeToString(sign)
		c.Proof = &Proof{
			Creator:            signer,
			Type:               ES256,
			SignatureValue:     signatureValue,
			VerificationMethod: method,
		}
		return nil
	default:
		return errors.New("unsupported type")
	}
}
