package tack

import "testing"
import "time"

func TestEntryMemUse(t *testing.T) {
	e := entry{
		time.Now().UnixNano(),
		"foo",
	}
	expected := 8 + 3
	got := e.getMem()
	if got != expected {
		t.Errorf("Expected %v to equal %v", expected, got)
	}

	expected = 3
	got = e.getValueMem()
	if got != expected {
		t.Errorf("Expected %v to equal %v", expected, got)
	}
}
