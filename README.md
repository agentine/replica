# replica

A modern, generics-based deep copy library for Go. Drop-in replacement for
[`mitchellh/copystructure`](https://github.com/mitchellh/copystructure).

## Features

- **Type-safe generics API** — no type assertions needed
- **Cycle detection** — returns error instead of stack overflow
- **Struct tags** — `copy:"ignore"` and `copy:"shallow"`
- **Functional options** — per-call configuration, no global state
- **Zero dependencies** — pure Go, no external packages
- **Compatibility layer** — `CopyAny`/`MustAny` for gradual migration

## Install

```
go get github.com/agentine/replica
```

## Usage

```go
import "github.com/agentine/replica"

// Type-safe deep copy
original := MyStruct{Name: "hello", Tags: []string{"a", "b"}}
copied, err := replica.Copy(original)

// Panic on error (for variable initialization)
copied := replica.Must(original)

// With options
copied, err := replica.CopyWith(original,
    replica.WithShallowTypes(reflect.TypeOf(time.Time{})),
)
```

## Struct Tags

```go
type Example struct {
    Normal  string
    Skip    string `copy:"ignore"`   // field is zeroed in copy
    Shallow *Big   `copy:"shallow"`  // shallow copy only
}
```

## Migration from copystructure

| Before (copystructure) | After (replica) |
|---|---|
| `copystructure.Copy(v)` | `replica.CopyAny(v)` |
| `copystructure.Must(v)` | `replica.MustAny(v)` |
| Type assertion: `result.(MyType)` | `replica.Copy(v)` (no assertion needed) |
| `copystructure.Copiers[t] = fn` | `replica.WithCopier(t, fn)` |
| `copystructure.ShallowCopiers[t] = true` | `replica.WithShallowTypes(t)` |

## License

MIT
