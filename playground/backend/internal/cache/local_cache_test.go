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
	pb "beam.apache.org/playground/backend/internal/api"
	"context"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestLocalCache_Get(t *testing.T) {
	preparedKey := GetKey(uuid.New(), Tag_StatusTag)

	preparedMap := make(map[string]Item)
	preparedMap[preparedKey] = Item{
		value:      pb.Status_STATUS_ERROR,
		created:    time.Time{},
		expiration: time.Now().Add(time.Minute).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[string]Item
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "LocalCache_GetExistKey",
			fields: fields{
				cleanupInterval: 5 * time.Second,
				items:           preparedMap,
			},
			args: args{
				key: preparedKey,
			},
			want:    pb.Status_STATUS_ERROR,
			wantErr: false,
		},
		{
			name: "LocalCache_GetNotExistKey",
			fields: fields{
				cleanupInterval: 5 * time.Second,
				items:           make(map[string]Item),
			},
			args: args{
				key: preparedKey,
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
			got, err := ls.Get(tt.args.key)
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
	type fields struct {
		cleanupInterval time.Duration
		items           map[string]Item
	}
	type args struct {
		key     string
		value   interface{}
		expTime time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "LocalCache_Set",
			fields: fields{
				cleanupInterval: 5 * time.Second,
				items:           make(map[string]Item),
			},
			args: args{
				key:     GetKey(uuid.New(), Tag_RunOutput),
				value:   "TEST_VALUE",
				expTime: 5 * time.Second,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			ls.Set(tt.args.key, tt.args.value, tt.args.expTime)
			_, err := ls.Get(tt.args.key)
			if err != nil {
				t.Errorf("Key %v with %v not set in cache", tt.args.key, tt.args.value)
			}
		})
	}
}

func TestLocalCache_clearItems(t *testing.T) {
	preparedMap := make(map[string]Item)
	keys := []string{"key1", "key2"}
	for _, key := range keys {
		preparedMap[key] = Item{
			value:      pb.Status_STATUS_ERROR,
			created:    time.Time{},
			expiration: time.Now().Add(time.Minute).UnixNano(),
		}
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[string]Item
	}
	type args struct {
		keys []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "LocalCache_clearItems",
			fields: fields{
				cleanupInterval: 5 * time.Second,
				items:           preparedMap,
			},
			args: args{keys: keys},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			ls.clearItems(tt.args.keys)
		})
	}
}

func TestLocalCache_expiredKeys(t *testing.T) {
	preparedMap := make(map[string]Item)
	preparedMap["expiredKey"] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Nanosecond).UnixNano(),
	}
	preparedMap["notExpiredKey"] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Minute).UnixNano(),
	}
	type fields struct {
		cleanupInterval time.Duration
		items           map[string]Item
	}
	tests := []struct {
		name     string
		fields   fields
		wantKeys []string
	}{
		{
			name: "LocalCache_expiredKeys",
			fields: fields{
				cleanupInterval: 2 * time.Nanosecond,
				items:           preparedMap,
			},
			wantKeys: []string{"expiredKey"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			time.Sleep(2 * time.Nanosecond)
			if gotKeys := ls.expiredKeys(); !reflect.DeepEqual(gotKeys, tt.wantKeys) {
				t.Errorf("expiredKeys() = %v, want %v", gotKeys, tt.wantKeys)
			}
		})
	}
}

func TestLocalCache_startGC(t *testing.T) {
	preparedMap := make(map[string]Item)
	preparedMap["expiredKey"] = Item{
		value:      "TEST_VALUE",
		created:    time.Time{},
		expiration: time.Now().Add(time.Nanosecond).UnixNano(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type fields struct {
		cleanupInterval time.Duration
		items           map[string]Item
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "LocalCache_startGC",
			fields: fields{
				cleanupInterval: 2 * time.Nanosecond,
				items:           preparedMap,
			},
			args: args{ctx: ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LocalCache{
				cleanupInterval: tt.fields.cleanupInterval,
				items:           tt.fields.items,
			}
			go ls.startGC(tt.args.ctx)
		})
	}
}
