package urlquery

import (
	"net/url"
	"testing"
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
	Kas
	Bmm
}

func TestUnmarshal(t *testing.T) {
	k := &KasA{}
	v := url.Values{}
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

	var asd any
	asd = k

	if err := Parse(v, asd); err != nil {
		t.Error(err)
	}
}
