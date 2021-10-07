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
	"github.com/google/uuid"
	"os"
)

type Tag string

const (
	TagStatus        Tag = "STATUS_TAG"
	TagRunOutput     Tag = "RUN_OUTPUT_TAG"
	TagCompileOutput Tag = "COMPILE_OUTPUT_TAG"
)

type client interface {
	// Get returns value from storage by pipelineId and tag.
	// If storage contains value by pipelineId and tag returns (value, true, nil).
	// If storage doesn't contain value by pipelineId and tag returns ("", false, nil).
	// If some error occurs method returns ("", false, err).
	get(key string) (interface{}, bool, error)

	// SetOrUpdate adds value to storage by pipelineId and tag.
	setOrUpdate(key string, value interface{})
}

type Storage struct {
	client client
}

// GetNewStorage returns new Storage to save and read value
func GetNewStorage() *Storage {
	switch os.Getenv("storage") {
	default:
		return &Storage{client: NewLocalStorageClient()}
	}
}

func (s *Storage) Get(pipelineId uuid.UUID, tag Tag) (interface{}, bool, error) {
	key := getKey(pipelineId, tag)
	return s.client.get(key)
}

func (s *Storage) Set(pipelineId uuid.UUID, tag Tag, value interface{}) {
	key := getKey(pipelineId, tag)
	s.client.setOrUpdate(key, value)
}

// getKey returns key for storage by id and tag
func getKey(pipelineId uuid.UUID, tag Tag) string {
	return fmt.Sprintf("%s_%s", pipelineId, tag)
}
