package gostree

import "testing"

func Test_foo(t *testing.T) {
	result := foo()
	expected := "bar"

	if result != expected {
		t.Errorf("expected %s, but got %s", expected, result)
	}
}
