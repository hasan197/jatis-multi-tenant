package di

import (
	"sync"
)

// Container adalah struktur yang menampung semua dependencies
type Container struct {
	mu    sync.RWMutex
	deps  map[string]interface{}
	close []func() error
}

// NewContainer membuat instance baru dari Container
func NewContainer() *Container {
	return &Container{
		deps:  make(map[string]interface{}),
		close: make([]func() error, 0),
	}
}

// Register mendaftarkan dependency ke container
func (c *Container) Register(name string, dep interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deps[name] = dep
}

// Get mengambil dependency dari container
func (c *Container) Get(name string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dep, ok := c.deps[name]
	return dep, ok
}

// RegisterCloser mendaftarkan fungsi cleanup untuk dependency
func (c *Container) RegisterCloser(fn func() error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.close = append(c.close, fn)
}

// Close menutup semua dependencies yang terdaftar
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	var errs []error
	for _, closeFn := range c.close {
		if err := closeFn(); err != nil {
			errs = append(errs, err)
		}
	}
	
	if len(errs) > 0 {
		return errs[0] // Return first error for simplicity
	}
	return nil
} 