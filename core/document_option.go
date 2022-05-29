package core

type Option interface {
	apply(doc *Document)
}

type authentication struct {
	method string
}

func (a *authentication) apply(doc *Document) {
	doc.Authentication = append(doc.Authentication, a.method)
}

func WithAuthentication(method string) Option {
	return &authentication{
		method: method,
	}
}
