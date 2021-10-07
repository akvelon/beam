// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"fmt"
	"sync"
	"time"
)

const (
	maxCacheSizeElements = 100
	defaultExpiration    = 10 * time.Second
	cleanupInterval      = 5 * time.Second
)

type LocalStorageClient struct {
	sync.RWMutex
	defaultExpiration    time.Duration
	cleanupInterval      time.Duration
	maxCacheSizeElements int
	items                map[string]Item
	queueHead            string
	queueTail            string
}

type Item struct {
	Value      interface{}
	Created    time.Time
	Next       string
	Previous   string
	Expiration int64
}

// NewLocalStorageClient creates new instance of local storage client
func NewLocalStorageClient() *LocalStorageClient {
	items := make(map[string]Item)
	ls := LocalStorageClient{
		defaultExpiration:    defaultExpiration,
		cleanupInterval:      cleanupInterval,
		items:                items,
		maxCacheSizeElements: maxCacheSizeElements,
	}

	go ls.startGC()

	return &ls
}

func (ls *LocalStorageClient) setOrUpdate(key string, value interface{}) {
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

func (ls *LocalStorageClient) get(key string) (interface{}, bool, error) {

	ls.RLock()
	item, found := ls.items[key]
	ls.RUnlock()
	if !found {
		return "", found, fmt.Errorf("item with key %s not found", key)
	}

	if item.Expiration > 0 {

		if time.Now().UnixNano() > item.Expiration {
			ls.delete(key)
			return "", false, fmt.Errorf("item with key %s is expired", key)
		}
	}

	ls.removeItemFromQueue(key)
	ls.setupNewHead(key)

	return item.Value, found, nil

}

func (ls *LocalStorageClient) removeItemFromQueue(key string) {
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

func (ls *LocalStorageClient) setupNewHead(newHeadKey string) {
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

func (ls *LocalStorageClient) delete(key string) {
	ls.removeItemFromQueue(key)
	ls.Lock()
	delete(ls.items, key)
	ls.Unlock()

}

func (ls *LocalStorageClient) startGC() {
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

func (ls *LocalStorageClient) expiredKeys() (keys []string) {
	ls.RLock()

	defer ls.RUnlock()

	for key, item := range ls.items {
		if time.Now().UnixNano() > item.Expiration && item.Expiration > 0 {
			keys = append(keys, key)
		}
	}
	return
}

func (ls *LocalStorageClient) clearItems(keys []string) {

	for _, key := range keys {
		ls.delete(key)
	}
}
