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
	"reflect"
	"testing"
	"time"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		want    *Storage
		wantErr bool
	}{
		{
			name: "NewStorage",
			want: &Storage{
				&LocalStorageClient{
					defaultExpiration:    defaultExpiration,
					cleanupInterval:      cleanupInterval,
					items:                make(map[string]Item),
					maxCacheSizeElements: maxCacheSizeElements,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStorage()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Get(t *testing.T) {
	pipelineId := uuid.New()
	value := "TEST_VALUE"
	preparedKey := fmt.Sprintf("%s_%s", pipelineId, Tag_StatusTag)
	preparedMap := make(map[string]Item)
	preparedMap[preparedKey] = Item{
		value:      value,
		created:    time.Now(),
		expiration: time.Now().Add(time.Minute).UnixNano(),
	}
	type fields struct {
		client client
	}
	type args struct {
		pipelineId uuid.UUID
		tag        Tag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "key doesn't exist",
			fields: fields{
				&LocalStorageClient{
					defaultExpiration:    defaultExpiration,
					cleanupInterval:      cleanupInterval,
					items:                make(map[string]Item),
					maxCacheSizeElements: maxCacheSizeElements,
				},
			},
			args: args{
				pipelineId: pipelineId,
				tag:        Tag_StatusTag,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "key exists",
			fields: fields{
				&LocalStorageClient{
					defaultExpiration:    defaultExpiration,
					cleanupInterval:      cleanupInterval,
					items:                preparedMap,
					maxCacheSizeElements: maxCacheSizeElements,
				},
			},
			args: args{
				pipelineId: pipelineId,
				tag:        Tag_StatusTag,
			},
			want:    value,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				client: tt.fields.client,
			}
			got, err := s.Get(tt.args.pipelineId, tt.args.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Set(t *testing.T) {
	type fields struct {
		client client
	}
	type args struct {
		pipelineId uuid.UUID
		tag        Tag
		value      interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "set",
			fields: fields{
				&LocalStorageClient{
					defaultExpiration:    defaultExpiration,
					cleanupInterval:      cleanupInterval,
					items:                make(map[string]Item),
					maxCacheSizeElements: maxCacheSizeElements,
				},
			},
			args: args{
				pipelineId: uuid.New(),
				tag:        Tag_StatusTag,
				value:      "TEST_VALUE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				client: tt.fields.client,
			}
			s.Set(tt.args.pipelineId, tt.args.tag, tt.args.value)
		})
	}
}

func Test_getKey(t *testing.T) {
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
			name: "getKey",
			args: args{
				pipelineId: pipelineId,
				tag:        Tag_StatusTag,
			},
			want: fmt.Sprintf("%s_%s", pipelineId, Tag_StatusTag),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getKey(tt.args.pipelineId, tt.args.tag); got != tt.want {
				t.Errorf("getKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
