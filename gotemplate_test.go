package gotemplate

import "testing"

func TestRemoveMe(t *testing.T) {
	if RemoveMe("hello", "world") != "helloworld" {
		t.Error("I expected helloworld")
	}
}
