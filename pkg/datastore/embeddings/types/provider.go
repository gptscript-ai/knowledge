package types

import cg "github.com/philippgille/chromem-go"

type EmbeddingModelProvider interface {
	Name() string
	EmbeddingFunc() (cg.EmbeddingFunc, error)
	Config() any
}
