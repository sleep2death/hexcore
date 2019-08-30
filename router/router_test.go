package router

import (
	"testing"

	"github.com/sleep2death/hexcore/actions"
	"gopkg.in/go-playground/assert.v1"
)

func TestRouter(t *testing.T) {
	r := New()
	r.Handle(Card, "/strike/:card_id/:target_id", func(ps Params) actions.Action {
		cid := ps.ByName("card_id")
		tid := ps.ByName("target_id")
		// t.Logf("card_id: %s, target_id: %s", ps.ByName("card_id"), ps.ByName("target_id"))
		assert.Equal(t, "abc", cid)
		assert.Equal(t, "def", tid)
		return nil
	})

	r.Serve(Card, "strike/abc/def")
	r.Serve(Card, "Strike/abc/def")
	r.Serve(Card, "Strike/abc/def/")
	r.Serve(Card, "strike/abc/def/")
}
