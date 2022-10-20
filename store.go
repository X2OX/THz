package THz

func (c *Context) Get(key any) (value any, exists bool) {
	c.mux.RLock()
	value, exists = c.keys[key]
	c.mux.RUnlock()
	return
}

func (c *Context) Set(key, val any) {
	c.mux.Lock()
	c.keys[key] = val
	c.mux.Unlock()
}

func (c *Context) Delete(key any) {
	c.mux.Lock()
	delete(c.keys, key)
	c.mux.Unlock()
}

func (c *Context) Reset() {
	c.mux.Lock()
	c.keys = make(map[any]any)
	c.mux.Unlock()
}

type Store[K, V any] struct {
	ctx *Context
}

func (s Store[K, V]) Get(key K) V {
	if val, exists := s.ctx.Get(key); exists {
		if v, ok := val.(V); ok {
			return v
		}
	}
	return *new(V)
}

func (s Store[K, V]) Set(k K, v V)              { s.ctx.Set(k, v) }
func (s Store[K, V]) Delete(k K)                { s.ctx.Delete(k) }
func GetStore[K, V any](c *Context) Store[K, V] { return Store[K, V]{ctx: c} }
