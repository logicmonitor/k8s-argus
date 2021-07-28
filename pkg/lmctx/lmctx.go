package lmctx

import "sync"

// LMContext will be used to pass metadata from one to other, for ex: log context.
type LMContext struct {
	kv sync.Map
}

// Extract returns value against key.
func (lc *LMContext) Extract(key string) interface{} {
	val, _ := lc.kv.Load(key)
	return val
}

// Set will store key val in map
func (lc *LMContext) Set(key string, val interface{}) {
	lc.kv.Store(key, val)
}

// Copy copy
func (lc *LMContext) Copy() *LMContext {
	// copying variable instead of pointer to make deep copy
	l2 := NewLMContext()
	lc.kv.Range(func(key, value interface{}) bool {
		l2.Set(key.(string), value)
		return true
	})
	return l2
}

// NewLMContext creates new context object
func NewLMContext() *LMContext {
	ctx := &LMContext{
		kv: sync.Map{},
	}

	return ctx
}

// LMContextWith creates new context object
func (lc *LMContext) LMContextWith(m map[string]interface{}) *LMContext {
	newLctx := lc.Copy()
	for k, v := range m {
		newLctx.Set(k, v)
	}
	return newLctx
}
