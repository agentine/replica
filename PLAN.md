# Replica — Modern Deep Copy Library for Go

## Overview

**Replaces:** `github.com/mitchellh/copystructure` (1,558 importers, 357 stars, archived July 2024)

**Package:** `github.com/agentine/replica`

A generics-era deep copy library for Go. Creates complete, independent copies of arbitrary Go values including structs, maps, slices, pointers, and nested combinations. Eliminates the archived `mitchellh/reflectwalk` dependency by inlining optimized traversal logic.

## Why Replace copystructure

- **Archived:** Repository archived July 2024, last release April 2022. Single maintainer (mitchellh) moved on from Go.
- **Pre-generics API:** Uses `interface{}` throughout, requiring type assertions on every call.
- **Global mutable state:** `Copiers` and `ShallowCopiers` are package-level maps — not safe for concurrent configuration.
- **Archived dependency:** Depends on `mitchellh/reflectwalk` (also archived, 1,503 importers).
- **No cycle detection:** Circular references cause stack overflow.
- **Fragmented alternatives:** tiendc/go-deepcopy (167 stars), brunoga/deep (110 stars), jinzhu/copier (different use case) — none have significant traction.

## Architecture

### Core API

```go
// Type-safe deep copy using generics
func Copy[T any](v T) (T, error)

// Panics on error — for variable initialization
func Must[T any](v T) T

// Configurable copy with functional options
func CopyWith[T any](v T, opts ...Option) (T, error)
```

### Options (functional options pattern)

```go
type Option func(*config)

func WithShallowTypes(types ...reflect.Type) Option  // Shallow copy specific types
func WithCopier(t reflect.Type, fn CopierFunc) Option // Custom copier for type
func WithLocking(lock bool) Option                     // Lock sync.Locker fields
```

### CopierFunc

```go
type CopierFunc func(any) (any, error)
```

### Struct Tags

```go
type Example struct {
    Normal  string
    Skip    string `copy:"ignore"`   // Skip this field
    Shallow *Big   `copy:"shallow"`  // Shallow copy only
}
```

## Major Components

1. **Core copier** (`copy.go`): Generic `Copy[T]`, `Must[T]`, `CopyWith[T]` functions. Dispatches by reflect.Kind.
2. **Struct walker** (`walk.go`): Inline struct/map/slice traversal. Replaces reflectwalk dependency. Handles unexported fields (skip), embedded structs, and struct tags.
3. **Cycle detector** (`cycle.go`): Pointer-based cycle detection using a visited set. Returns error on circular reference instead of stack overflow.
4. **Options** (`options.go`): Functional options for per-call configuration. No global mutable state.
5. **Compatibility** (`compat.go`): Optional `interface{}`-based wrappers for gradual migration from copystructure.

## Supported Types

- Primitives (bool, int*, uint*, float*, complex*, string)
- Pointers (with cycle detection)
- Structs (exported fields, struct tags)
- Maps (deep copy keys and values)
- Slices and arrays
- Interfaces
- sync.Locker types (optional locking during copy)
- time.Time and other common value types (shallow by default)
- Custom types via CopierFunc

## Compatibility Layer

```go
// Drop-in replacement for copystructure.Copy
func CopyAny(v any) (any, error)

// Drop-in replacement for copystructure.Must
func MustAny(v any) any
```

## Deliverables

- `github.com/agentine/replica` Go module
- Zero external dependencies
- Full test suite with cycle detection, struct tags, edge cases
- Benchmarks vs copystructure
- README with migration guide from copystructure
- Go 1.21+ (generics + `any` alias)

## Non-Goals

- Copying between different struct types (use jinzhu/copier for that)
- Serialization/deserialization
- Copying unexported struct fields (skip with warning, matching copystructure behavior)
