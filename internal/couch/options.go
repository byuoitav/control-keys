package couch

import "github.com/go-kivik/couchdb/v3"

type options struct {
	authFunc   interface{}
	uiConfigDB string
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(o *options) {
	f(o)
}

func WithBasicAuth(username, password string) Option {
	return optionFunc(func(o *options) {
		o.authFunc = couchdb.BasicAuth(username, password)
	})
}
