package core

import (
	"encoding/json"
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

	Authentication     *VerificationMethod   `json:"authentication"`
	VerificationMethod []*VerificationMethod `json:"verificationMethod"`
	Service            *Service              `json:"service"`
}

func NewDocument(id string, options ...Option) *Document {
	doc := &Document{
		ID: id,
	}
	for _, echo := range options {
		echo.apply(doc)
	}
	return doc
}

func (d *Document) PutVerifyMethod(method *VerificationMethod) {
	d.VerificationMethod = append(d.VerificationMethod, method)
}

func (d *Document) GetVerifyMethod(id string) *VerificationMethod {
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

type Service struct {
	ID              string   `json:"id"`
	Type            []string `json:"type"`
	ServiceEndpoint string   `json:"serviceEndpoint"`
}
