package storage

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type LocalStorage struct {
	sync.RWMutex
	defaultExpiration    time.Duration
	cleanupInterval      time.Duration
	maxCacheSizeElements int
	items                map[string]Item
	queueHead            string
	queueTail            string
}

type Item struct {
	Value      string
	Created    time.Time
	Next       string
	Previous   string
	Expiration int64
}

func NewLocalStorage(defaultExpiration, cleanupInterval time.Duration, maxCacheSizeElements int) (*LocalStorage, error) {
	if maxCacheSizeElements <= 0 {
		return nil, errors.New("trying to create cache with 0 elements")
	}
	items := make(map[string]Item)
	ls := LocalStorage{
		defaultExpiration:    defaultExpiration,
		cleanupInterval:      cleanupInterval,
		items:                items,
		maxCacheSizeElements: maxCacheSizeElements,
	}

	if cleanupInterval > 0 {
		go ls.startGC()
	}
	return &ls, nil

}

func (ls *LocalStorage) SetOrUpdate(key, value string) {
	expiration := time.Now().Add(ls.defaultExpiration).UnixNano()

	ls.RLock()
	item, found := ls.items[key]
	ls.RUnlock()
	if found {
		ls.removeItemFromQueue(key)

	} else {
		item = Item{}
	}

	item.Value = value
	item.Expiration = expiration
	item.Created = time.Now()
	ls.Lock()
	ls.items[key] = item

	if ls.queueTail == "" {
		ls.queueTail = key
	}
	ls.Unlock()
	ls.setupNewHead(key)

	if len(ls.items) <= ls.maxCacheSizeElements {
		return
	}
	ls.delete(ls.queueTail)
}

func (ls *LocalStorage) Get(key string) (string, error) {

	ls.RLock()
	item, found := ls.items[key]
	ls.RUnlock()
	if !found {
		return "", fmt.Errorf("item with key %s not found", key)
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			ls.delete(key)
			return "", fmt.Errorf("item with key %s is expired", key)
		}
	}

	ls.removeItemFromQueue(key)
	ls.setupNewHead(key)

	return item.Value, nil

}

func (ls *LocalStorage) removeItemFromQueue(key string) {
	ls.Lock()
	defer ls.Unlock()
	item, found := ls.items[key]
	if !found {
		return
	}
	if item.Previous == "" { // head
		ls.queueHead = item.Next
	} else {
		previous := ls.items[item.Previous]
		previous.Next = item.Next
		ls.items[item.Previous] = previous
	}
	if item.Next == "" { //tail
		ls.queueTail = item.Previous
	} else {
		next := ls.items[item.Next]
		next.Previous = item.Previous
		ls.items[item.Next] = next
	}
}

func (ls *LocalStorage) setupNewHead(newHeadKey string) {
	ls.Lock()
	item := ls.items[newHeadKey]
	if ls.queueHead != "" {
		currentHead := ls.items[ls.queueHead]
		currentHead.Previous = newHeadKey
		ls.items[ls.queueHead] = currentHead
	}
	item.Next = ls.queueHead
	item.Previous = ""
	ls.items[newHeadKey] = item
	ls.queueHead = newHeadKey
	ls.Unlock()
}

func (ls *LocalStorage) delete(key string) {
	ls.removeItemFromQueue(key)
	ls.Lock()
	delete(ls.items, key)
	ls.Unlock()

}

func (ls *LocalStorage) startGC() {
	for {
		<-time.After(ls.cleanupInterval)

		if ls.items == nil {
			return
		}

		if keys := ls.expiredKeys(); len(keys) != 0 {
			ls.clearItems(keys)
		}
	}
}

func (ls *LocalStorage) expiredKeys() (keys []string) {
	ls.RLock()

	defer ls.RUnlock()

	for key, item := range ls.items {
		if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
			keys = append(keys, key)
		}
	}
	return
}

func (ls *LocalStorage) clearItems(keys []string) {

	for _, key := range keys {
		ls.delete(key)
	}
}
