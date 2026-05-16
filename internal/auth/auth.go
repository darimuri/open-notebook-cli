package auth

import (
	"net/http"
)

type Middleware struct {
	apiKey string
}

func NewMiddleware(apiKey string) *Middleware {
	return &Middleware{apiKey: apiKey}
}

func (m *Middleware) AddAuth(req *http.Request) {
	if m.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+m.apiKey)
	}
}