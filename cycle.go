package replica

import (
	"errors"
	"reflect"
)

// ErrCycleDetected is returned when a circular reference is found during copy.
var ErrCycleDetected = errors.New("replica: circular reference detected")

// visited tracks seen pointers to detect cycles during deep copy.
type visited struct {
	ptrs map[uintptr]bool
}

func newVisited() *visited {
	return &visited{ptrs: make(map[uintptr]bool)}
}

// seen checks if a pointer-like value has been visited. Returns true if it was
// already seen (cycle detected). Non-pointer kinds are never cycles.
func (v *visited) seen(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		if val.IsNil() {
			return false
		}
		ptr := val.Pointer()
		if v.ptrs[ptr] {
			return true
		}
		v.ptrs[ptr] = true
	}
	return false
}

// remove deletes a pointer from the visited set (for backtracking).
func (v *visited) remove(val reflect.Value) {
	switch val.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		if !val.IsNil() {
			delete(v.ptrs, val.Pointer())
		}
	}
}
