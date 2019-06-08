package github

import (
	"golang.org/x/oauth2"
)

// TokenSource authorizes GitHub access and provides a name for your GitHub
// user or organization
type TokenSource interface {
	Name() string

	oauth2.TokenSource
}

type staticSource struct {
	oauth2.TokenSource
	name string
}

// NewStaticTokenSource returns a new token source that just uses a static token
func NewStaticTokenSource(name string, t oauth2.Token) TokenSource {
	return &staticSource{
		TokenSource: oauth2.StaticTokenSource(&t),
		name:        name,
	}
}

func (s *staticSource) Name() string { return s.name }
