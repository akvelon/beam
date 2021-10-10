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
	"github.com/google/uuid"
	"time"
)

type Tag string

const (
	Tag_StatusTag        Tag = "STATUS_TAG"
	Tag_RunOutputTag     Tag = "RUN_OUTPUT_TAG"
	Tag_CompileOutputTag Tag = "COMPILE_OUTPUT_TAG"
)

type Cache interface {
	// Get returns value from cache by key.
	Get(pipelineId uuid.UUID, tag Tag) (interface{}, error)

	// Set adds value to cache by key.
	Set(pipelineId uuid.UUID, tag Tag, value interface{}, expTime time.Duration)
}

// GetNewCache returns new cache to save and read value
func GetNewCache(cacheType string) Cache {
	switch cacheType {
	default:
		return newLocalCache()
	}
}
