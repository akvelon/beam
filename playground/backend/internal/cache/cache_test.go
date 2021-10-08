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
	"reflect"
	"testing"
	"time"
)

func TestGetNewCache1(t *testing.T) {
	tests := []struct {
		name string
		want Cache
	}{
		{
			name: "NewCache",
			want: &LocalCache{
				cleanupInterval: 5 * time.Second,
				items:           make(map[string]Item),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNewCache(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNewCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	pipelineId := uuid.New()
	type args struct {
		pipelineId uuid.UUID
		tag        Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "GetKey",
			args: args{
				pipelineId: pipelineId,
				tag:        Tag_StatusTag,
			},
			want: fmt.Sprintf("%s_%s", pipelineId, Tag_StatusTag),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKey(tt.args.pipelineId, tt.args.tag); got != tt.want {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
