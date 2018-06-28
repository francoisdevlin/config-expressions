package confexpr

import (
	"testing"
)

func TestMatchStateString(t *testing.T) {
	expect(Missing.String(), "Missing", t)
	expect(Incomplete.String(), "Incomplete", t)
	expect(Collision.String(), "Collision", t)
	expect(Complete.String(), "Complete", t)
}
