package authorization

import (
	"bytes"
	"errors"

	thz "go.x2ox.com/THz"
	"go.x2ox.com/sorbifolia/jwt"
)

type Auth[T any] struct {
	j         *jwt.JWT[T]
	abort     bool
	abortFunc func(ctx *thz.Context)
	store     bool
	storeKey  string
}

func New[T any](j *jwt.JWT[T], abort bool, abortFunc func(ctx *thz.Context), store bool, storeKey string) *Auth[T] {
	g := &Auth[T]{j: j, abort: abort, abortFunc: abortFunc, store: store, storeKey: storeKey}
	if g.abort && g.abortFunc == nil {
		g.abortFunc = func(c *thz.Context) { c.Status(401).Abort() }
	}

	if g.store && g.storeKey == "" {
		g.storeKey = "JWT"
	}

	return g
}

func (g *Auth[T]) Middleware() thz.Handler {
	return func(c *thz.Context) {
		claims, err := g.Parse(c.Header("Authorization"))
		if err != nil {
			if g.abort {
				g.abortFunc(c)
			}
			return
		}

		if g.store {
			c.Set(g.storeKey, claims.Data)
		}

		c.Next()
	}
}

func (g *Auth[T]) Parse(authHeader []byte) (*jwt.Claims[T], error) {
	if len(authHeader) < 7 || !bytes.EqualFold(authHeader[:7], []byte("Bearer ")) {
		return nil, errors.New("invalid authorization")
	}

	return g.j.Parse(string(authHeader[7:]))
}
