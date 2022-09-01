package urlquery

import (
	"net/url"
	"testing"
	"time"
)

type Kas struct {
	A   string    `query:"a"`
	Arr [3]string `query:"arr"`
	As  []string  `query:"as"`
}

type Bmm struct {
	Ba any `query:"ba"`
	Bu any `query:"bu"`
}

type KasA struct {
	T time.Time `query:"t"`
	Kas
	Bmm
}

func TestUnmarshal(t *testing.T) {
	var (
		asd any
		k   = &KasA{}
		v   = url.Values{}
	)
	v.Set("t", time.Now().Format(time.RFC1123))
	v.Set("a", "asdasd")
	v.Add("arr", "1")
	v.Add("arr", "2")
	v.Add("as", "1")
	v.Add("as", "2")
	v.Add("as", "3")
	v.Add("as", "4")
	v.Add("as", "5")
	v.Add("ba", "bb")
	v.Add("bu", "bb")
	v.Add("ba", "bbb")
	v.Add("ba", "bbbb")

	asd = k

	if err := Parse(v, asd); err != nil {
		t.Error(err)
	}
}
