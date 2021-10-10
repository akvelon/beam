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
	"fmt"
	"github.com/google/uuid"
	"time"
)

const cleanupInterval = 5 * time.Second

type LocalCache struct {
	cleanupInterval time.Duration
	items           map[uuid.UUID]map[Tag]Item
}

type Item struct {
	value      interface{}
	created    time.Time
	expiration int64
}

func newLocalCache() *LocalCache {
	items := make(map[uuid.UUID]map[Tag]Item)
	ls := &LocalCache{
		cleanupInterval: cleanupInterval,
		items:           items,
	}

	go ls.startGC()
	return ls

}

func (ls *LocalCache) Get(pipelineId uuid.UUID, tag Tag) (interface{}, error) {
	item, found := ls.items[pipelineId][tag]
	if !found {
		return nil, fmt.Errorf("Item with pipelineId: %s and tag: %s not found.", pipelineId, tag)
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			delete(ls.items[pipelineId], tag)
			return nil, fmt.Errorf("item with pipelineId: %s and tag: %s is expired", pipelineId, tag)
		}
	}

	return item.value, nil
}

func (ls *LocalCache) Set(pipelineId uuid.UUID, tag Tag, value interface{}, expTime time.Duration) {
	expiration := time.Now().Add(expTime).UnixNano()
	if _, ok := ls.items[pipelineId]; !ok {
		ls.items[pipelineId] = make(map[Tag]Item)
	}
	item := ls.items[pipelineId][tag]
	item.value = value
	item.expiration = expiration
	item.created = time.Now()
	ls.items[pipelineId][tag] = item
}

func (ls *LocalCache) startGC() {
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

func (ls *LocalCache) expiredKeys() (keys map[uuid.UUID][]Tag) {
	keys = make(map[uuid.UUID][]Tag)
	for pipeline := range ls.items {
		for tag, item := range ls.items[pipeline] {
			if time.Now().UnixNano() > item.expiration && item.expiration > 0 {
				keys[pipeline] = append(keys[pipeline], tag)
			}
		}
	}
	return
}

func (ls *LocalCache) clearItems(keys map[uuid.UUID][]Tag) {
	for pipeline := range keys {
		for _, tag := range keys[pipeline] {
			delete(ls.items[pipeline], tag)
		}
		if len(ls.items[pipeline]) == 0 {
			delete(ls.items, pipeline)
		}
	}
}
