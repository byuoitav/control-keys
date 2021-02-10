package couch

import (
	"context"
	"fmt"

	"github.com/go-kivik/kivik/v3"
)

type dataService struct {
	client     *kivik.Client
	uiConfigDB string
}

// New creates a new ConfigService, created a couchdb client pointed at url.
func New(ctx context.Context, url string, opts ...Option) (*dataService, error) {
	client, err := kivik.New("couch", url)
	if err != nil {
		return nil, fmt.Errorf("unable to build client: %w", err)
	}

	return NewWithClient(ctx, client, opts...)
}

// NewWithClient creates a new ConfigService using the given client.
func NewWithClient(ctx context.Context, client *kivik.Client, opts ...Option) (*dataService, error) {
	options := options{
		uiConfigDB: "ui-configuration",
	}

	for _, o := range opts {
		o.apply(&options)
	}

	if options.authFunc != nil {
		if err := client.Authenticate(ctx, options.authFunc); err != nil {
			return nil, fmt.Errorf("unable to authenticate: %w", err)
		}
	}

	return &dataService{
		client:     client,
		uiConfigDB: options.uiConfigDB,
	}, nil
}
