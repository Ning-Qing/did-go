package core

import (
	"encoding/base64"
	"log"
)

func ValidateClaim(claim *Claim, doc *Document) bool {
	method := doc.GetVerifMethod(claim.Proof.VerificationMethod)
	if method == nil {
		return false
	}
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
