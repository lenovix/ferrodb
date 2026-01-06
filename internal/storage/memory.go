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
	data []map[string]Item
	mu   sync.RWMutex
}

func NewMemoryStore(dbCount int, cleanupIntervalSec int) *MemoryStore {
	data := make([]map[string]Item, dbCount)
	for i := range data {
		data[i] = make(map[string]Item)
	}

	store := &MemoryStore{data: data}

	go store.cleanupExpiredKeys(time.Duration(cleanupIntervalSec) * time.Second)

	return store
}

func (m *MemoryStore) Set(db int, key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item := m.data[db][key]
	item.Value = value
	item.ExpireAt = 0
	m.data[db][key] = item
}

func (m *MemoryStore) Get(db int, key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[db][key]
	if !ok {
		return "", false
	}

	if item.ExpireAt > 0 && time.Now().Unix() > item.ExpireAt {
		delete(m.data[db], key)
		return "", false
	}

	return item.Value, true
}

func (m *MemoryStore) Del(db int, key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data[db], key)
}

func (m *MemoryStore) ExpireAt(db int, key string, timestamp int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[db][key]
	if !ok {
		return false
	}

	item.ExpireAt = timestamp
	m.data[db][key] = item
	return true
}

func (m *MemoryStore) TTL(db int, key string) int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[db][key]
	if !ok {
		return -2
	}

	if item.ExpireAt == 0 {
		return -1
	}

	now := time.Now().Unix()
	if now >= item.ExpireAt {
		delete(m.data[db], key)
		return -2
	}

	return item.ExpireAt - now
}

func (m *MemoryStore) Persist(db int, key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, ok := m.data[db][key]
	if !ok || item.ExpireAt == 0 {
		return false
	}

	item.ExpireAt = 0
	m.data[db][key] = item
	return true
}

func (m *MemoryStore) Snapshot() map[int]map[string]Item {
	m.mu.RLock()
	defer m.mu.RUnlock()

	snap := make(map[int]map[string]Item)
	now := time.Now().Unix()

	for db, kv := range m.data {
		dbSnap := make(map[string]Item)
		for k, v := range kv {
			if v.ExpireAt > 0 && now > v.ExpireAt {
				continue
			}
			dbSnap[k] = v
		}
		if len(dbSnap) > 0 {
			snap[db] = dbSnap
		}
	}

	return snap
}

func (m *MemoryStore) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := 0
	for _, kv := range m.data {
		total += len(kv)
	}
	return total
}

func (m *MemoryStore) cleanupExpiredKeys(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()

		m.mu.Lock()
		for db, kv := range m.data {
			for k, v := range kv {
				if v.ExpireAt > 0 && now > v.ExpireAt {
					delete(m.data[db], k)
				}
			}
		}
		m.mu.Unlock()
	}
}
