package example

import "testing"

type museumCtx struct {
	Alex, Blake, Casey suspect
}

type suspect struct {
	did_id, is_truthful bool
}

type museumAct struct {
	name   string
	mutate func(*museumCtx)
}

var museumInput = []museumAct{
	{name: "Alex", mutate: func(mc *museumCtx) { mc.Blake.did_id = true }},       // Blake did it
	{name: "Blake", mutate: func(mc *museumCtx) { mc.Casey.did_id = false }},     // Casey did NOT do it
	{name: "Casey", mutate: func(mc *museumCtx) { mc.Alex.is_truthful = false }}, // Alex is lying
}

// one man guilty, one truth teller, two lying
// who stole the painting?
func TestMuseumTheft(t *testing.T) {

}
