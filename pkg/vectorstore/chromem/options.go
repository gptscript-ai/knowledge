package chromem

import (
	"errors"
	"github.com/philippgille/chromem-go"
)

// Option is a function type that can be used to modify the client.
type Option func(p *Store)

// WithDB is an option for setting the database to use.
func WithDB(d *chromem.DB) Option {
	return func(p *Store) {
		p.db = d
	}
}

// WithCollection is an option for setting the collection to use.
func WithCollection(c *chromem.Collection) Option {
	return func(p *Store) {
		p.collection = c
	}
}

func WithEmbeddingFunc(f chromem.EmbeddingFunc) Option {
	return func(p *Store) {
		p.embeddingFunc = f
	}
}

func applyClientOptions(opts ...Option) (Store, error) {
	o := &Store{}

	for _, opt := range opts {
		opt(o)
	}

	if o.db == nil {
		return Store{}, errors.New("missing database")
	}

	if o.collection == nil {
		return Store{}, errors.New("missing collection")
	}

	if o.embeddingFunc == nil {
		return Store{}, errors.New("missing embedding function")
	}

	return *o, nil
}
