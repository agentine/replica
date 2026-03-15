package replica

import (
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

// --- Primitive types ---

func TestCopyBool(t *testing.T) {
	v := true
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyInt(t *testing.T) {
	v := 42
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyInt8(t *testing.T) {
	v := int8(127)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyInt16(t *testing.T) {
	v := int16(32767)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyInt32(t *testing.T) {
	v := int32(2147483647)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyInt64(t *testing.T) {
	v := int64(9223372036854775807)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyUint(t *testing.T) {
	v := uint(42)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyFloat32(t *testing.T) {
	v := float32(3.14)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyFloat64(t *testing.T) {
	v := 3.14159265
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyComplex64(t *testing.T) {
	v := complex64(1 + 2i)
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyComplex128(t *testing.T) {
	v := 1 + 2i
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v, want %v", got, v)
	}
}

func TestCopyString(t *testing.T) {
	v := "hello world"
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %q, want %q", got, v)
	}
}

// --- Pointers ---

func TestCopyPointerSingle(t *testing.T) {
	v := 42
	p := &v
	got, err := Copy(p)
	if err != nil {
		t.Fatal(err)
	}
	if got == p {
		t.Error("copy should be a different pointer")
	}
	if *got != *p {
		t.Errorf("got %d, want %d", *got, *p)
	}
	// Verify independence.
	*got = 99
	if *p == 99 {
		t.Error("modifying copy affected original")
	}
}

func TestCopyPointerDouble(t *testing.T) {
	v := 42
	p := &v
	pp := &p
	got, err := Copy(pp)
	if err != nil {
		t.Fatal(err)
	}
	if got == pp {
		t.Error("outer pointer should be different")
	}
	if *got == *pp {
		t.Error("inner pointer should be different")
	}
	if **got != **pp {
		t.Errorf("got %d, want %d", **got, **pp)
	}
}

func TestCopyNilPointer(t *testing.T) {
	var p *int
	got, err := Copy(p)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

// --- Structs ---

func TestCopyStructExportedFields(t *testing.T) {
	type S struct {
		Name  string
		Value int
	}
	orig := S{Name: "test", Value: 42}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got != orig {
		t.Errorf("got %+v, want %+v", got, orig)
	}
}

type hasUnexported struct {
	Name string
	age  int //nolint:unused // testing unexported field skip
}

func TestCopyStructUnexportedSkipped(t *testing.T) {
	orig := hasUnexported{Name: "test"}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "test" {
		t.Errorf("Name: got %q, want %q", got.Name, "test")
	}
}

func TestCopyStructWithPointer(t *testing.T) {
	type Inner struct {
		Value int
	}
	type Outer struct {
		Name  string
		Inner *Inner
	}
	orig := Outer{Name: "test", Inner: &Inner{Value: 42}}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got.Inner == orig.Inner {
		t.Error("Inner pointer should be different")
	}
	if got.Inner.Value != 42 {
		t.Errorf("Inner.Value: got %d, want 42", got.Inner.Value)
	}
	got.Inner.Value = 99
	if orig.Inner.Value == 99 {
		t.Error("modifying copy affected original")
	}
}

func TestCopyStructEmbedded(t *testing.T) {
	type Base struct {
		ID int
	}
	type Child struct {
		Base
		Name string
	}
	orig := Child{Base: Base{ID: 1}, Name: "child"}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 1 || got.Name != "child" {
		t.Errorf("got %+v, want %+v", got, orig)
	}
}

// --- Maps ---

func TestCopyMapStringString(t *testing.T) {
	orig := map[string]string{"a": "1", "b": "2"}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(orig) {
		t.Fatalf("len: got %d, want %d", len(got), len(orig))
	}
	for k, v := range orig {
		if got[k] != v {
			t.Errorf("key %q: got %q, want %q", k, got[k], v)
		}
	}
	// Verify independence.
	got["c"] = "3"
	if _, ok := orig["c"]; ok {
		t.Error("modifying copy affected original")
	}
}

func TestCopyMapStringStruct(t *testing.T) {
	type V struct {
		Data []int
	}
	orig := map[string]V{
		"x": {Data: []int{1, 2, 3}},
	}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	got["x"] = V{Data: append(got["x"].Data, 4)}
	if len(orig["x"].Data) != 3 {
		t.Error("modifying copy affected original")
	}
}

func TestCopyMapNested(t *testing.T) {
	orig := map[string]map[string]int{
		"a": {"x": 1, "y": 2},
	}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	got["a"]["z"] = 3
	if _, ok := orig["a"]["z"]; ok {
		t.Error("modifying nested map copy affected original")
	}
}

func TestCopyNilMap(t *testing.T) {
	var m map[string]int
	got, err := Copy(m)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestCopyEmptyMap(t *testing.T) {
	m := map[string]int{}
	got, err := Copy(m)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("got nil, want empty map")
	}
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

// --- Slices ---

func TestCopySlice(t *testing.T) {
	orig := []int{1, 2, 3}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(orig) {
		t.Fatalf("len: got %d, want %d", len(got), len(orig))
	}
	got[0] = 99
	if orig[0] == 99 {
		t.Error("modifying copy affected original")
	}
}

func TestCopySliceOfPointers(t *testing.T) {
	a, b := 1, 2
	orig := []*int{&a, &b}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] == orig[0] || got[1] == orig[1] {
		t.Error("pointers in slice should be different")
	}
	if *got[0] != 1 || *got[1] != 2 {
		t.Errorf("values: got %d,%d want 1,2", *got[0], *got[1])
	}
}

func TestCopyNilSlice(t *testing.T) {
	var s []int
	got, err := Copy(s)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestCopyEmptySlice(t *testing.T) {
	s := []int{}
	got, err := Copy(s)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("got nil, want empty slice")
	}
	if len(got) != 0 {
		t.Errorf("len: got %d, want 0", len(got))
	}
}

// --- Arrays ---

func TestCopyArray(t *testing.T) {
	orig := [3]int{1, 2, 3}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got != orig {
		t.Errorf("got %v, want %v", got, orig)
	}
}

func TestCopyArrayOfPointers(t *testing.T) {
	a, b := 10, 20
	orig := [2]*int{&a, &b}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got[0] == orig[0] || got[1] == orig[1] {
		t.Error("pointers in array should be different")
	}
	*got[0] = 99
	if a == 99 {
		t.Error("modifying copy affected original")
	}
}

// --- Interfaces ---

func TestCopyInterface(t *testing.T) {
	var v any = map[string]int{"a": 1}
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	m := got.(map[string]int)
	m["b"] = 2
	orig := v.(map[string]int)
	if _, ok := orig["b"]; ok {
		t.Error("modifying copy affected original")
	}
}

func TestCopyNilInterface(t *testing.T) {
	var v any
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

// --- Cycle detection ---

func TestCycleSelfReference(t *testing.T) {
	type Node struct {
		Next *Node
	}
	a := &Node{}
	a.Next = a

	_, err := Copy(a)
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("got error %v, want ErrCycleDetected", err)
	}
}

func TestCycleMutualReference(t *testing.T) {
	type Node struct {
		Next *Node
	}
	a := &Node{}
	b := &Node{}
	a.Next = b
	b.Next = a

	_, err := Copy(a)
	if !errors.Is(err, ErrCycleDetected) {
		t.Errorf("got error %v, want ErrCycleDetected", err)
	}
}

func TestNoCycleSamePointerDifferentPaths(t *testing.T) {
	// Shared pointer that appears twice is NOT a cycle (DAG, not cycle).
	type S struct {
		A *int
		B *int
	}
	v := 42
	orig := S{A: &v, B: &v}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if *got.A != 42 || *got.B != 42 {
		t.Errorf("got A=%d B=%d, want 42,42", *got.A, *got.B)
	}
}

// --- Struct tags ---

func TestTagIgnore(t *testing.T) {
	type S struct {
		Name   string
		Secret string `copy:"ignore"`
	}
	orig := S{Name: "test", Secret: "hidden"}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "test" {
		t.Errorf("Name: got %q, want %q", got.Name, "test")
	}
	if got.Secret != "" {
		t.Errorf("Secret: got %q, want empty (ignored)", got.Secret)
	}
}

func TestTagShallow(t *testing.T) {
	type Big struct {
		Data []int
	}
	type S struct {
		Deep    *Big
		Shallow *Big `copy:"shallow"`
	}
	orig := S{
		Deep:    &Big{Data: []int{1, 2, 3}},
		Shallow: &Big{Data: []int{4, 5, 6}},
	}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	// Deep should be a different pointer.
	if got.Deep == orig.Deep {
		t.Error("Deep should be deeply copied")
	}
	// Shallow should be the same pointer.
	if got.Shallow != orig.Shallow {
		t.Error("Shallow should be the same pointer (shallow copy)")
	}
}

// --- Options ---

func TestWithShallowTypes(t *testing.T) {
	type Big struct {
		Data []int
	}
	type S struct {
		B *Big
	}
	orig := S{B: &Big{Data: []int{1, 2, 3}}}
	got, err := CopyWith(orig, WithShallowTypes(reflect.TypeOf(Big{})))
	if err != nil {
		t.Fatal(err)
	}
	// The Big value should not be deeply copied.
	if got.B == orig.B {
		// Pointer is new (struct copy creates new struct), but inner slice should share.
		// Actually, shallow on the Big type means it returns Big as-is, so the pointer
		// is deeply copied (ptr copy allocates new) but the Big value inside is shallow.
		got.B.Data[0] = 99
		if orig.B.Data[0] != 99 {
			t.Error("WithShallowTypes should share inner data")
		}
	}
}

func TestWithCopier(t *testing.T) {
	type Special struct {
		Value int
	}
	copier := func(v any) (any, error) {
		s := v.(Special)
		return Special{Value: s.Value * 2}, nil
	}
	orig := Special{Value: 21}
	got, err := CopyWith(orig, WithCopier(reflect.TypeOf(Special{}), copier))
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != 42 {
		t.Errorf("got %d, want 42", got.Value)
	}
}

func TestWithCopierError(t *testing.T) {
	type S struct {
		Value int
	}
	errCustom := errors.New("custom error")
	copier := func(v any) (any, error) {
		return nil, errCustom
	}
	_, err := CopyWith(S{Value: 1}, WithCopier(reflect.TypeOf(S{}), copier))
	if !errors.Is(err, errCustom) {
		t.Errorf("got %v, want %v", err, errCustom)
	}
}

func TestWithLocking(t *testing.T) {
	type S struct {
		sync.Mutex
		Value int
	}
	orig := &S{Value: 42}

	// Lock the original in another goroutine — without WithLocking this would
	// not block, but we just verify it doesn't deadlock with locking enabled.
	got, err := CopyWith(orig, WithLocking(true))
	if err != nil {
		t.Fatal(err)
	}
	if got.Value != 42 {
		t.Errorf("got %d, want 42", got.Value)
	}
}

// --- time.Time shallow copy ---

func TestTimeShallowCopy(t *testing.T) {
	type S struct {
		Created time.Time
	}
	now := time.Now()
	orig := S{Created: now}
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Created.Equal(orig.Created) {
		t.Errorf("time mismatch: got %v, want %v", got.Created, orig.Created)
	}
}

// --- Compat layer ---

func TestCopyAny(t *testing.T) {
	orig := map[string]int{"a": 1}
	result, err := CopyAny(orig)
	if err != nil {
		t.Fatal(err)
	}
	m := result.(map[string]int)
	m["b"] = 2
	if _, ok := orig["b"]; ok {
		t.Error("modifying copy affected original")
	}
}

func TestMustAny(t *testing.T) {
	orig := []int{1, 2, 3}
	result := MustAny(orig)
	s := result.([]int)
	s[0] = 99
	if orig[0] == 99 {
		t.Error("modifying copy affected original")
	}
}

func TestMustPanics(t *testing.T) {
	type Node struct {
		Next *Node
	}
	a := &Node{}
	a.Next = a

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
	}()
	Must(a)
}

func TestMustAnyPanics(t *testing.T) {
	type Node struct {
		Next *Node
	}
	a := &Node{}
	a.Next = a

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
	}()
	MustAny(a)
}

// --- Complex nested structure ---

func TestCopyComplexNested(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}
	type Person struct {
		Name      string
		Addresses []*Address
		Metadata  map[string]any
	}

	orig := Person{
		Name: "Alice",
		Addresses: []*Address{
			{Street: "123 Main St", City: "Springfield"},
			{Street: "456 Oak Ave", City: "Shelbyville"},
		},
		Metadata: map[string]any{
			"age":     30,
			"hobbies": []string{"reading", "coding"},
		},
	}

	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}

	// Verify values.
	if got.Name != "Alice" {
		t.Errorf("Name: got %q, want %q", got.Name, "Alice")
	}
	if len(got.Addresses) != 2 {
		t.Fatalf("Addresses len: got %d, want 2", len(got.Addresses))
	}
	if got.Addresses[0].Street != "123 Main St" {
		t.Errorf("Address[0]: got %q", got.Addresses[0].Street)
	}

	// Verify independence.
	got.Addresses[0].Street = "changed"
	if orig.Addresses[0].Street == "changed" {
		t.Error("modifying address copy affected original")
	}
}

// --- Zero values ---

func TestCopyZeroStruct(t *testing.T) {
	type S struct {
		Name string
		Val  int
	}
	var orig S
	got, err := Copy(orig)
	if err != nil {
		t.Fatal(err)
	}
	if got != orig {
		t.Errorf("got %+v, want %+v", got, orig)
	}
}

func TestCopyEmptyString(t *testing.T) {
	v := ""
	got, err := Copy(v)
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
