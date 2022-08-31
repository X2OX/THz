package authorization

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	"go.x2ox.com/THz"
	"go.x2ox.com/sorbifolia/jwt"
)

const testKey = "5d94fce0d33407b4d085023dcbf79c04e98d9f7d56ee87e981a07b02a704" +
	"ab34c7a04587614b1aa0fbba6e76446c99ccea11696dcc72d26ac84f478902ec4627"

type A struct {
	User string `json:"user"`
}

func newJWT() {
	bts, err := hex.DecodeString(testKey)
	if err != nil {
		panic(err)
	}
	priKey := ed25519.PrivateKey(bts)
	_jwt := jwt.New(jwt.EdDSA, priKey, priKey.Public(), jwt.Claims[A]{})
	g := New(_jwt, true,
		func(c *THz.Context) {
			c.Abort().Status(http.StatusOK).Text("fail")
		},
		true, "user",
	)

	fr := &fasthttp.Request{Header: fasthttp.RequestHeader{}}
	fr.Header.Set("Authorization", "Bearer "+_jwt.MustGenerate(
		*jwt.NewClaims(&A{User: "k"}).SetExpiresAt(time.Now().Add(time.Hour * 6)),
	))

	c, e := g.Parse(fr)
	fmt.Println(e)
	fmt.Println(c.Data)
}

func TestA(t *testing.T) {
	newJWT()
}
