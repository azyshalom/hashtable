package hashtable

import (
    "sync"
    "strings"
)

type Bucket map[string]interface{}
type DoFunc func(interface{}, interface{})
type HashFunc func(string) int

type Hash struct {
    buckets  []Bucket
    size     int
    count    int
    mutex    sync.RWMutex
    hashFunc HashFunc
}

func NewHash(size int, f HashFunc) *Hash {
    h := new(Hash)
    h.mutex = sync.RWMutex{}
    h.size = size
    h.count = 0
    h.buckets = make([]Bucket, size, size)
    h.hashFunc = f

    if h.hashFunc == nil {
         h.hashFunc = h.hash
    }

    for i := 0; i < h.size; i++ {
        h.buckets[i] = make(map[string]interface{})
    }

    return h
}

func (h *Hash) Count() int {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    return h.count
}

func (h *Hash) Add(key string, value interface{}) {
    h.mutex.Lock()
    defer h.mutex.Unlock()
    idx := h.hashFunc(key)
    h.buckets[idx][key] = value
    h.count++
}

func (h *Hash) Get(key string) interface{} {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    idx := h.hashFunc(key)
    if _, ok := h.buckets[idx][key]; ok {
        return h.buckets[idx][key]
    }

    return nil
}

func (h *Hash) Delete(key string) { 
    h.mutex.Lock()
    defer h.mutex.Unlock()
    idx := h.hashFunc(key)
    if _, ok := h.buckets[idx][key]; ok {
        delete(h.buckets[idx], key)
        h.count--
    }
}

func (h *Hash) DoAll(f DoFunc, ctx interface{}) {
    h.mutex.RLock()
    defer h.mutex.RUnlock()
    for i := 0; i < h.size ; i++ {
        for _, value := range h.buckets[i] {
            f(ctx, value)
        }
    }
}

func (h *Hash) hash(key string) int {
    sum := 0
    str := strings.ToUpper(key)
    for i := 0; i < len(str); i++ {
        sum+= int(str[i])
    }
    return sum % h.size
}  
