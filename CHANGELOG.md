# Changelog

## 0.1.0 (2026-03-16)

Initial release — generics-era deep copy library for Go, replacing `mitchellh/copystructure`.

### Features

- Generic `Copy[T]`, `Must[T]`, and `CopyWith[T]` functions with full `reflect.Kind` dispatch
- Cycle detection via pointer visited set (returns error instead of stack overflow)
- Struct tags: `copy:"ignore"` and `copy:"shallow"`
- Functional options: `WithShallowTypes`, `WithCopier`, `WithLocking`
- `time.Time` shallow-copied by default
- `sync.Locker` locking support for concurrent-safe structs
- `CopyAny` / `MustAny` compatibility layer for `copystructure` migration
- Zero external dependencies
- 94.7% test coverage, 48 tests, 4 benchmarks
- CI tested on Go 1.23 and 1.24
