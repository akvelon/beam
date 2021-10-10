//Licensed to the Apache Software Foundation (ASF) under one or more
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
	pb "beam.apache.org/playground/backend/internal/api"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestLocalCache_Get(t *testing.T) {
	preparedId, _ := uuid.NewUUID()
	preparedTag := Tag_CompileOutputTag
	value := "TEST_VALUE"
	preparedMap := make(map[uuid.UUID]map[Tag]Item)
	preparedMap[preparedId] = make(map[Tag]Item)
	preparedMap[preparedId][preparedTag] = Item{
		value:      value,
		created:    time.Time{},
		expiration: time.Now().Add(time.Minute).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[uuid.UUID]map[Tag]Item
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
			name: "LocalCache_GetExistItem",
			fields: fields{
				cleanupInterval: cleanupInterval,
				items:           preparedMap,
			},
			args: args{
				pipelineId: preparedId,
				tag:        preparedTag,
			},
			want:    value,
			wantErr: false,
		},
		{
			name: "LocalCache_GetNotExistItem",
			fields: fields{
				cleanupInterval: cleanupInterval,
				items:           make(map[uuid.UUID]map[Tag]Item),
			},
			args: args{
				pipelineId: preparedId,
				tag:        preparedTag,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			got, err := ls.Get(tt.args.pipelineId, tt.args.tag)
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

func TestLocalCache_Set(t *testing.T) {
	preparedId, _ := uuid.NewUUID()
	type fields struct {
		cleanupInterval time.Duration
		items           map[uuid.UUID]map[Tag]Item
	}
	type args struct {
		pipelineId uuid.UUID
		tag        Tag
		value      interface{}
		expTime    time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "LocalCache_Set",
			fields: fields{
				cleanupInterval: cleanupInterval,
				items:           make(map[uuid.UUID]map[Tag]Item),
			},
			args: args{
				pipelineId: preparedId,
				tag:        Tag_RunOutputTag,
				value:      "TEST_VALUE",
				expTime:    time.Minute,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			ls.Set(tt.args.pipelineId, tt.args.tag, tt.args.value, tt.args.expTime)
			_, err := ls.Get(tt.args.pipelineId, tt.args.tag)
			if err != nil {
				t.Errorf("Item with pipelineId: %s and tag: %v not set in cache.", tt.args.pipelineId, tt.args.tag)
			}
		})
	}
}

func TestLocalCache_startGC(t *testing.T) {
	preparedId, _ := uuid.NewUUID()
	preparedTag := Tag_CompileOutputTag
	value := "TEST_VALUE"
	preparedMap := make(map[uuid.UUID]map[Tag]Item)
	preparedMap[preparedId] = make(map[Tag]Item)
	preparedMap[preparedId][preparedTag] = Item{
		value:      value,
		created:    time.Time{},
		expiration: time.Now().Add(time.Microsecond).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[uuid.UUID]map[Tag]Item
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "LocalCache_startGC",
			fields: fields{
				cleanupInterval: time.Microsecond,
				items:           preparedMap,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			go ls.startGC()
			time.Sleep(time.Millisecond)
			if len(tt.fields.items) != 0 {
				t.Errorf("Item with pipelineId: %s and tag: %v not deleted in time.", preparedId, preparedTag)
			}
		})
	}
}

func TestLocalCache_expiredKeys(t *testing.T) {
	preparedId1, _ := uuid.NewUUID()
	preparedId2, _ := uuid.NewUUID()
	expiredValue := "EXPIRED_VALUE"
	notExpiredValue := pb.Status_STATUS_FINISHED
	preparedMap := make(map[uuid.UUID]map[Tag]Item)
	preparedMap[preparedId1] = make(map[Tag]Item)
	preparedMap[preparedId2] = make(map[Tag]Item)
	preparedMap[preparedId1][Tag_StatusTag] = Item{
		value:      expiredValue,
		created:    time.Time{},
		expiration: time.Now().Add(time.Microsecond).UnixNano(),
	}
	preparedMap[preparedId1][Tag_RunOutputTag] = Item{
		value:      notExpiredValue,
		created:    time.Time{},
		expiration: time.Now().Add(time.Second).UnixNano(),
	}
	preparedMap[preparedId2][Tag_StatusTag] = Item{
		value:      expiredValue,
		created:    time.Time{},
		expiration: time.Now().Add(time.Microsecond).UnixNano(),
	}
	preparedMap[preparedId2][Tag_RunOutputTag] = Item{
		value:      expiredValue,
		created:    time.Time{},
		expiration: time.Now().Add(time.Microsecond).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[uuid.UUID]map[Tag]Item
	}
	tests := []struct {
		name     string
		fields   fields
		wantKeys map[uuid.UUID][]Tag
	}{
		{
			name: "LocalCache_expiredKeys",
			fields: fields{
				cleanupInterval: cleanupInterval,
				items:           preparedMap,
			},
			wantKeys: map[uuid.UUID][]Tag{preparedId1: {Tag_StatusTag}, preparedId2: {Tag_StatusTag, Tag_RunOutputTag}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			if gotKeys := ls.expiredKeys(); !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("expiredKeys() = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func TestLocalCache_clearItems(t *testing.T) {
	preparedId1, _ := uuid.NewUUID()
	preparedId2, _ := uuid.NewUUID()
	preparedMap := make(map[uuid.UUID]map[Tag]Item)
	preparedMap[preparedId1] = make(map[Tag]Item)
	preparedMap[preparedId2] = make(map[Tag]Item)
	preparedMap[preparedId1][Tag_RunOutputTag] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Second).UnixNano(),
	}
	preparedMap[preparedId2][Tag_RunOutputTag] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Second).UnixNano(),
	}
	preparedMap[preparedId2][Tag_CompileOutputTag] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Second).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[uuid.UUID]map[Tag]Item
	}
	type args struct {
		keys map[uuid.UUID][]Tag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "LocalCache_clearItems",
			fields: fields{
				cleanupInterval: cleanupInterval,
				items:           preparedMap,
			},
			args: args{keys: map[uuid.UUID][]Tag{preparedId1: {Tag_RunOutputTag}, preparedId2: {Tag_CompileOutputTag}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			ls.clearItems(tt.args.keys)
			if _, err := ls.Get(preparedId2, Tag_RunOutputTag); err != nil {
				t.Error(err)
			}
			if _, err := ls.Get(preparedId1, Tag_RunOutputTag); err == nil {
				t.Errorf("The desired item with pipelineId: %s and tag:%v has not been deleted.", preparedId2, Tag_RunOutputTag)
			}
			if _, found := tt.fields.items[preparedId1]; found {
				t.Errorf("The empty map which key: %s without item has not been deleted.", preparedId1)
			}
		})
	}
}
