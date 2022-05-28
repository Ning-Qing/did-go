package core

import (
	"crypto"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
)

type Controller struct {
	provider Provider
}

func NewController(provider Provider) *Controller {
	return &Controller{
		provider: provider,
	}
}

func (c *Controller) NewDID(method string, jwk Jwk) string {
	pubKey := jwk.GeneratePubKey()

	did := fmt.Sprintf("did:%s:%s", method, generate(pubKey))
	auth := fmt.Sprintf("%s#auth", did)

	doc := NewDocument(did, WithAuthentication(NewVerificationMethod(auth, did, pubKey)))

	c.provider.PutDocument(doc)
	return did
}

func (c *Controller) ValidateClaim(claim *Claim) bool {
	// doc := c.provider.GetDocument(claim.Issuer)
	return false
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
