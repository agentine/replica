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
