package core

type Option interface {
	apply(doc *Document)
}

type authentication struct {
	method *VerificationMethod
}

func (a *authentication) apply(doc *Document) {
	doc.Authentication = a.method
}

func WithAuthentication(method *VerificationMethod) Option {
	return &authentication{
		method: method,
	}
}
