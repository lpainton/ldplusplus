package game

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type fakePRNG struct {
	deflt int
	queue []int
}

func (f *fakePRNG) Intn(_ int) int {
	if len(f.queue) == 0 {
		return f.deflt
	}
	first, rest := f.queue[0], f.queue[1:]
	f.queue = rest
	return first
}

func TestNew(t *testing.T) {
	type testcase struct {
		name string
		Rules
	}

	tests := []testcase{
		{
			name: "New plumbs through the rules applied to it",
			Rules: Rules{
				Dice:  5,
				Wilds: [6]bool{true, false, false, false, false, false},
			},
		},
	}

	for _, test := range tests {
		g := New(test.Rules)
		if !cmp.Equal(g.Rules, test.Rules) {
			t.Errorf("test: %q failed...want: %v, got %v", test.name, test.Rules, g.Rules)
		}
	}
}
