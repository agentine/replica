package replica

import (
	"reflect"
)

// CopierFunc is a custom copy function for a specific type.
type CopierFunc func(any) (any, error)

// Option configures the behavior of CopyWith.
type Option func(*config)

type config struct {
	shallowTypes map[reflect.Type]bool
	copiers      map[reflect.Type]CopierFunc
	locking      bool
}

func newConfig(opts []Option) *config {
	c := &config{}
	for _, o := range opts {
		o(c)
	}
	return c
}

func (c *config) isShallow(t reflect.Type) bool {
	if c.shallowTypes == nil {
		return false
	}
	return c.shallowTypes[t]
}

func (c *config) getCopier(t reflect.Type) (CopierFunc, bool) {
	if c.copiers == nil {
		return nil, false
	}
	fn, ok := c.copiers[t]
	return fn, ok
}

// WithShallowTypes marks the given types for shallow copy only (pointer copy,
// no recursion into fields).
func WithShallowTypes(types ...reflect.Type) Option {
	return func(c *config) {
		if c.shallowTypes == nil {
			c.shallowTypes = make(map[reflect.Type]bool)
		}
		for _, t := range types {
			c.shallowTypes[t] = true
		}
	}
}

// WithCopier registers a custom copy function for the given type.
func WithCopier(t reflect.Type, fn CopierFunc) Option {
	return func(c *config) {
		if c.copiers == nil {
			c.copiers = make(map[reflect.Type]CopierFunc)
		}
		c.copiers[t] = fn
	}
}

// WithLocking enables locking of sync.Locker fields during copy.
func WithLocking(lock bool) Option {
	return func(c *config) {
		c.locking = lock
	}
}
