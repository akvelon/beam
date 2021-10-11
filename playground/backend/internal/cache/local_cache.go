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
	"sync"
	"time"
)

const cleanupInterval = 5 * time.Second

type LocalCache struct {
	sync.RWMutex
	cleanupInterval     time.Duration
	items               map[uuid.UUID]map[SubKey]interface{}
	pipelinesExpiration map[uuid.UUID]time.Time
}

func newLocalCache() *LocalCache {
	items := make(map[uuid.UUID]map[SubKey]interface{})
	pipelinesExpiration := make(map[uuid.UUID]time.Time)
	ls := &LocalCache{
		cleanupInterval:     cleanupInterval,
		items:               items,
		pipelinesExpiration: pipelinesExpiration,
	}

	go ls.startGC()
	return ls

}

func (lc *LocalCache) GetValue(pipelineId uuid.UUID, subKey SubKey) (interface{}, error) {
	lc.RLock()
	value, found := lc.items[pipelineId][subKey]
	if !found {
		return nil, fmt.Errorf("value with pipelineId: %s and subKey: %s not found", pipelineId, subKey)
	}
	expTime, found := lc.pipelinesExpiration[pipelineId]
	if !found {
		return nil, fmt.Errorf("expiration time of the pipeline: %s not found", pipelineId)
	}
	lc.RUnlock()

	if expTime.Before(time.Now()) {
		lc.Lock()
		delete(lc.items[pipelineId], subKey)
		lc.Unlock()
		return nil, fmt.Errorf("value with pipelineId: %s and subKey: %s is expired", pipelineId, subKey)
	}

	return value, nil
}

func (lc *LocalCache) SetValue(pipelineId uuid.UUID, subKey SubKey, value interface{}) {
	lc.Lock()
	defer lc.Unlock()

	_, ok := lc.items[pipelineId]
	if !ok {
		lc.items[pipelineId] = make(map[SubKey]interface{})
	}
	lc.items[pipelineId][subKey] = value
}

func (lc *LocalCache) SetExpTime(pipelineId uuid.UUID, expTime time.Duration) {
	lc.Lock()
	defer lc.Unlock()
	lc.pipelinesExpiration[pipelineId] = time.Now().Add(expTime)
}

func (lc *LocalCache) startGC() {
	for {
		<-time.After(lc.cleanupInterval)

		if lc.items == nil {
			return
		}

		if pipelines := lc.expiredPipelines(); len(pipelines) != 0 {
			lc.clearItems(pipelines)
		}
	}
}

func (lc *LocalCache) expiredPipelines() (pipelines []uuid.UUID) {
	lc.RLock()
	defer lc.RUnlock()
	for pipelineId, expTime := range lc.pipelinesExpiration {
		if expTime.Before(time.Now()) {
			pipelines = append(pipelines, pipelineId)
		}
	}
	return
}

func (lc *LocalCache) clearItems(pipelines []uuid.UUID) {
	for _, pipeline := range pipelines {
		lc.Lock()
		delete(lc.items, pipeline)
		lc.Unlock()
	}
}
