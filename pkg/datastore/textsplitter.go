package datastore

import lcgosplitter "github.com/tmc/langchaingo/textsplitter"

var (
	defaultLcgoSplitter = lcgosplitter.NewTokenSplitter(lcgosplitter.WithChunkSize(defaultChunkSize), lcgosplitter.WithChunkOverlap(defaultChunkOverlap), lcgosplitter.WithModelName(defaultTokenModel), lcgosplitter.WithEncodingName(defaultTokenEncoding))
)
