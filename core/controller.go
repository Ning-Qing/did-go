package core

import (
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"log"
)

const (
	SCHEME     = "did"
	METHODNAME = "example"
	BASEAUTH   = "baseauth"
)

type Controller struct {
	provider Provider
}

func NewController(provider Provider) *Controller {
	return &Controller{
		provider: provider,
	}
}

func (c *Controller) NewDID(jwk Jwk) string {
	pubKey := jwk.GeneratePubKey()
	url := &URL{
		Scheme:   SCHEME,
		Method:   METHODNAME,
		Path:     generate(pubKey),
		Fragment: BASEAUTH,
	}
	baseauth := NewVerificationMethod(url.DIDServer(), url.DID(), pubKey)
	doc := NewDocument(url.DID(), WithAuthentication(url.DIDServer()))
	doc.PutVerifyMethod(baseauth)
	c.provider.PutDocument(doc)
	return url.DID()
}

func (c *Controller) GetDocument(id string) *Document {
	return c.provider.GetDocument(id)
}

func (c *Controller) AuthenticationClaim(claim *Claim) bool {
	doc := c.provider.GetDocument(claim.Issuer)
	method := doc.GetAuthentication(claim.Proof.VerificationMethod)
	if method == nil {
		return false
	}
	return validateClaim(claim, method)
}

func generate(pubKey crypto.PublicKey) string {
	key := NewPubKeyFactory(pubKey)
	data, err := key.Marshal()
	if err != nil {
		log.Println(err)
	}
	hash := sha256.New()
	hash.Write(data)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func validateClaim(claim *Claim, method *VerificationMethod) bool {
	hash, err := claim.Hash()
	if err != nil {
		log.Println(err)
		return false
	}
	sign, err := base64.StdEncoding.DecodeString(claim.Proof.SignatureValue)
	if err != nil {
		log.Println(err)
		return false
	}
	return method.Verify(hash, sign)
}
