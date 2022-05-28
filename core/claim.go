package core

import (
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
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
	*Doc
	Proof *Proof `json:"proof"`
}

type Proof struct {
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

func (c *Claim) Sign(method string, privkey crypto.PrivateKey) error {
	key := NewPrivKeyFactory(privkey)
	digest, err := c.Hash()
	if err != nil {
		return err
	}
	sign, err := key.Sign(digest)
	if err != nil {
		return err
	}
	signatureValue := base64.StdEncoding.EncodeToString(sign)
	c.Proof = &Proof{
		Type:               RSA256,
		SignatureValue:     signatureValue,
		VerificationMethod: method,
	}
	return nil
}
