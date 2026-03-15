package replica

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

// defaultShallowTypes are always shallow-copied regardless of options.
var defaultShallowTypes = map[reflect.Type]bool{
	reflect.TypeOf(time.Time{}): true,
}

var syncLockerType = reflect.TypeOf((*sync.Locker)(nil)).Elem()
var syncMutexType = reflect.TypeOf(sync.Mutex{})
var syncRWMutexType = reflect.TypeOf(sync.RWMutex{})

// deepCopy recursively copies val using cfg and vis for cycle detection.
func deepCopy(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	if !val.IsValid() {
		return val, nil
	}

	t := val.Type()

	// Check custom copier first.
	if fn, ok := cfg.getCopier(t); ok {
		result, err := fn(val.Interface())
		if err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(result), nil
	}

	// Shallow types: return as-is (value copy).
	if defaultShallowTypes[t] || cfg.isShallow(t) {
		return val, nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		return copyPtr(val, cfg, vis)
	case reflect.Struct:
		return copyStruct(val, cfg, vis)
	case reflect.Map:
		return copyMap(val, cfg, vis)
	case reflect.Slice:
		return copySlice(val, cfg, vis)
	case reflect.Array:
		return copyArray(val, cfg, vis)
	case reflect.Interface:
		return copyInterface(val, cfg, vis)
	default:
		// Primitives (bool, int*, uint*, float*, complex*, string, chan, func, etc.)
		return val, nil
	}
}

func copyPtr(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	if val.IsNil() {
		return reflect.Zero(val.Type()), nil
	}

	if vis.seen(val) {
		return reflect.Value{}, ErrCycleDetected
	}
	defer vis.remove(val)

	elem, err := deepCopy(val.Elem(), cfg, vis)
	if err != nil {
		return reflect.Value{}, err
	}

	ptr := reflect.New(val.Type().Elem())
	ptr.Elem().Set(elem)
	return ptr, nil
}

func copyStruct(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	t := val.Type()
	result := reflect.New(t).Elem()

	// If locking is enabled and the struct implements sync.Locker, lock it.
	// Skip sync primitives themselves — they are the lock, not data to protect.
	if cfg.locking && val.CanAddr() && t != syncMutexType && t != syncRWMutexType && val.Addr().Type().Implements(syncLockerType) {
		locker := val.Addr().Interface().(sync.Locker)
		locker.Lock()
		defer locker.Unlock()
	}

	for i := range t.NumField() {
		field := t.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		// Parse struct tag.
		tag := parseTag(field)
		if tag == tagIgnore {
			continue
		}

		srcField := val.Field(i)
		if tag == tagShallow {
			result.Field(i).Set(srcField)
			continue
		}

		copied, err := deepCopy(srcField, cfg, vis)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("field %s: %w", field.Name, err)
		}
		result.Field(i).Set(copied)
	}

	return result, nil
}

func copyMap(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	if val.IsNil() {
		return reflect.Zero(val.Type()), nil
	}

	if vis.seen(val) {
		return reflect.Value{}, ErrCycleDetected
	}
	defer vis.remove(val)

	result := reflect.MakeMapWithSize(val.Type(), val.Len())

	iter := val.MapRange()
	for iter.Next() {
		copiedKey, err := deepCopy(iter.Key(), cfg, vis)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("map key: %w", err)
		}
		copiedVal, err := deepCopy(iter.Value(), cfg, vis)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("map value: %w", err)
		}
		result.SetMapIndex(copiedKey, copiedVal)
	}

	return result, nil
}

func copySlice(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	if val.IsNil() {
		return reflect.Zero(val.Type()), nil
	}

	if vis.seen(val) {
		return reflect.Value{}, ErrCycleDetected
	}
	defer vis.remove(val)

	result := reflect.MakeSlice(val.Type(), val.Len(), val.Cap())

	for i := range val.Len() {
		copied, err := deepCopy(val.Index(i), cfg, vis)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("slice[%d]: %w", i, err)
		}
		result.Index(i).Set(copied)
	}

	return result, nil
}

func copyArray(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	result := reflect.New(val.Type()).Elem()

	for i := range val.Len() {
		copied, err := deepCopy(val.Index(i), cfg, vis)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("array[%d]: %w", i, err)
		}
		result.Index(i).Set(copied)
	}

	return result, nil
}

func copyInterface(val reflect.Value, cfg *config, vis *visited) (reflect.Value, error) {
	if val.IsNil() {
		return reflect.Zero(val.Type()), nil
	}

	copied, err := deepCopy(val.Elem(), cfg, vis)
	if err != nil {
		return reflect.Value{}, err
	}

	return copied, nil
}

type tagAction int

const (
	tagNone    tagAction = iota
	tagIgnore            // copy:"ignore"
	tagShallow           // copy:"shallow"
)

func parseTag(field reflect.StructField) tagAction {
	tag := field.Tag.Get("copy")
	switch tag {
	case "ignore":
		return tagIgnore
	case "shallow":
		return tagShallow
	default:
		return tagNone
	}
}
