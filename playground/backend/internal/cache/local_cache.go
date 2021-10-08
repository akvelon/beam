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

package cache

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type LocalCache struct {
	cleanupInterval time.Duration
	items           map[string]Item
}

type Item struct {
	value      interface{}
	created    time.Time
	expiration int64
}

func newLocalCache() *LocalCache {
	items := make(map[string]Item)
	ls := &LocalCache{
		cleanupInterval: 5 * time.Second,
		items:           items,
	}
	ctx, cancel := context.WithCancel(context.Background())
	runtime.SetFinalizer(ls, cancel)
	go ls.startGC(ctx)
	return ls

}

func (ls *LocalCache) Get(key string) (interface{}, error) {
	item, found := ls.items[key]
	if !found {
		return nil, fmt.Errorf("item with key %s not found", key)
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			delete(ls.items, key)
			return nil, fmt.Errorf("item with key %s is expired", key)
		}
	}

	return item.value, nil
}

func (ls *LocalCache) Set(key string, value interface{}, expTime time.Duration) {
	expiration := time.Now().Add(expTime).UnixNano()

	item := ls.items[key]
	item.value = value
	item.expiration = expiration
	item.created = time.Now()
	ls.items[key] = item
}

func (ls *LocalCache) startGC(ctx context.Context) {
	ticker := time.NewTicker(ls.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			if keys := ls.expiredKeys(); len(keys) != 0 {
				ls.clearItems(keys)
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}

	}
}

func (ls *LocalCache) expiredKeys() (keys []string) {

	for key, item := range ls.items {
		if time.Now().UnixNano() > item.expiration && item.expiration > 0 {
			keys = append(keys, key)
		}
	}
	return
}

func (ls *LocalCache) clearItems(keys []string) {
	for _, key := range keys {
		delete(ls.items, key)
	}
}
