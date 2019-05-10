package chord

import (
	"errors"
)

type kvStore struct {
	storage map[string]string
}

func NewKVStore() *kvStore {
	return &kvStore{storage: make(map[string]string)}
}

func (kv *kvStore) Get(key string) (string, error) {
	val, ok := kv.storage[key]
	if ok {
		return val, nil
	}
	return "", errors.New("key not existing")
}

func (kv *kvStore) Put(key string, val string) error {
	kv.storage[key] = val
	return nil
}

func (kv *kvStore) Delete(key string) (string, error) {
	val := kv.storage[key]
	delete(kv.storage, key)
	return val, nil
}

func (kv *kvStore) KeyCount() int {
	return len(kv.storage)
}

// kv.kvMu.Lock()
// defer kv.kvMu.Unlock()
