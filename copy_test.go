package replica

import (
	"testing"
)

func TestCopyBasicStruct(t *testing.T) {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Name  string
		Tags  []string
		Inner *Inner
	}

	original := Outer{
		Name:  "test",
		Tags:  []string{"a", "b"},
		Inner: &Inner{Value: 42},
	}

	copied, err := Copy(original)
	if err != nil {
		t.Fatal(err)
	}

	// Verify values match.
	if copied.Name != original.Name {
		t.Errorf("Name: got %q, want %q", copied.Name, original.Name)
	}
	if len(copied.Tags) != len(original.Tags) {
		t.Fatalf("Tags len: got %d, want %d", len(copied.Tags), len(original.Tags))
	}
	if copied.Inner.Value != 42 {
		t.Errorf("Inner.Value: got %d, want 42", copied.Inner.Value)
	}

	// Verify independence.
	copied.Tags[0] = "changed"
	if original.Tags[0] == "changed" {
		t.Error("modifying copied slice affected original")
	}
	copied.Inner.Value = 99
	if original.Inner.Value == 99 {
		t.Error("modifying copied pointer affected original")
	}
}

func TestCopyCycleDetection(t *testing.T) {
	type Node struct {
		Next *Node
	}
	a := &Node{}
	a.Next = a // self-reference

	_, err := Copy(a)
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}
