package storage

import (
	"sync"
	"time"
)

type Item struct {
	Value    string
	ExpireAt int64 // unix timestamp (seconds), 0 = no TTL
}

type MemoryStore struct {
	data map[string]Item
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]Item),
	}

	go store.cleanupExpiredKeys()
	return store
}

func (m *MemoryStore) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := m.data[key]
	item.Value = value
	m.data[key] = item
}

func (m *MemoryStore) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[key]
	if !ok {
		return "", false
	}

	if item.ExpireAt > 0 && time.Now().Unix() > item.ExpireAt {
		delete(m.data, key)
		return "", false
	}

	return item.Value, true
}

func (m *MemoryStore) Del(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

func (m *MemoryStore) Expire(key string, seconds int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[key]
	if !ok {
		return false
	}

	item.ExpireAt = time.Now().Unix() + seconds
	m.data[key] = item
	return true
}

func (m *MemoryStore) cleanupExpiredKeys() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()

		m.mu.Lock()
		for k, v := range m.data {
			if v.ExpireAt > 0 && now > v.ExpireAt {
				delete(m.data, k)
			}
		}
		m.mu.Unlock()
	}
}

func (m *MemoryStore) ExpireAt(key string, timestamp int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[key]
	if !ok {
		return false
	}

	item.ExpireAt = timestamp
	m.data[key] = item
	return true
}

func (m *MemoryStore) Snapshot() map[string]Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snap := make(map[string]Item, len(m.data))
	now := time.Now().Unix()

	for k, v := range m.data {
		if v.ExpireAt > 0 && now > v.ExpireAt {
			continue
		}
		snap[k] = v
	}

	return snap
}
