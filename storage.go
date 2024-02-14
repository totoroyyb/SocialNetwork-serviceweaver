package main

import (
	"context"
	"sync"
)

type Hashtable interface {
  Put(ctx context.Context, key, value string) error
  Get(ctx context.Context, key string) (string, error)
}

// the hashtable needs to be manually managed so that it has only one instance and no replicas.
// Later on, for distributed sharded hastable, we can have more replicas as shards.
type hashtable struct {
  mu sync.Mutex
  data map[string]string
}

func (ht *hashtable) Put(_ context.Context, key, value string) error {
  ht.mu.Lock()
  defer ht.mu.Unlock()
  ht.data[key] = value
  return nil
}

func (ht *hashtable) Get(_ context.Context, key string) (string, error) {
  ht.mu.Lock()
  defer ht.mu.Unlock()
  // test, err := ht.data[key]
  return ht.data[key], nil
}
