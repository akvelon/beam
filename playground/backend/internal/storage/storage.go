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
	Tag_StatusTag     Tag = "STATUS_TAG"
	Tag_RunOutput     Tag = "RUN_OUTPUT_TAG"
	Tag_CompileOutput Tag = "COMPILE_OUTPUT_TAG"
)

type client interface {
	// get returns value from storage by key.
	get(key string) (interface{}, error)

	// set adds value to storage by key.
	set(key string, value interface{})
}

type Storage struct {
	client client
}

// NewStorage returns new Storage to save and read value
func NewStorage() (*Storage, error) {
	switch os.Getenv("storage") {
	default:
		client, err := NewLocalStorageClient()
		if err != nil {
			return nil, err
		}
		return &Storage{client}, nil
	}
}

func (s *Storage) Get(pipelineId uuid.UUID, tag Tag) (interface{}, error) {
	return s.client.get(getKey(pipelineId, tag))
}

func (s *Storage) Set(pipelineId uuid.UUID, tag Tag, value interface{}) {
	s.client.set(getKey(pipelineId, tag), value)
}

// getKey returns key for storage by id and tag
func getKey(pipelineId uuid.UUID, tag Tag) string {
	return fmt.Sprintf("%s_%s", pipelineId, tag)
}
