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
package main

import (
	pb "beam.apache.org/playground/backend/internal/api/v1"
	"beam.apache.org/playground/backend/internal/cache"
	"beam.apache.org/playground/backend/internal/cache/local"
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"reflect"
	"strings"
	"testing"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener
var cacheService cache.Cache

func TestMain(m *testing.M) {
	server := setup()
	defer teardown(server)
	m.Run()
}

func setup() *grpc.Server {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	cacheService = local.New(context.Background())

	pb.RegisterPlaygroundServiceServer(s, &playgroundController{cacheService: cacheService})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	return s
}

func teardown(server *grpc.Server) {
	server.Stop()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestPlaygroundController_RunCode(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlaygroundServiceClient(conn)
	code := pb.RunCodeRequest{
		Code: "test",
		Sdk:  pb.Sdk_SDK_JAVA,
	}
	pipelineMeta, err := client.RunCode(ctx, &code)
	if err != nil {
		t.Fatalf("runCode failed: %v", err)
	}
	log.Printf("Response: %+v", pipelineMeta)
}

func TestPlaygroundController_CheckStatus(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlaygroundServiceClient(conn)
	pipelineMeta := pb.CheckStatusRequest{
		PipelineUuid: uuid.NewString(),
	}
	status, err := client.CheckStatus(ctx, &pipelineMeta)
	if err != nil {
		t.Fatalf("runCode failed: %v", err)
	}
	log.Printf("Response: %+v", status)
}

func TestPlaygroundController_GetCompileOutput(t *testing.T) {
	ctx := context.Background()
	pipelineId := uuid.New()
	compileOutput := "MOCK_COMPILE_OUTPUT"
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlaygroundServiceClient(conn)

	type args struct {
		ctx  context.Context
		info *pb.GetCompileOutputRequest
	}
	tests := []struct {
		name    string
		prepare func()
		args    args
		want    *pb.GetRunOutputResponse
		wantErr bool
	}{
		{
			name:    "pipelineId doesn't exist",
			prepare: func() {},
			args: args{
				ctx:  ctx,
				info: &pb.GetCompileOutputRequest{PipelineUuid: pipelineId.String()},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "run output exist",
			prepare: func() {
				_ = cacheService.SetValue(ctx, pipelineId, cache.CompileOutput, compileOutput)
			},
			args: args{
				ctx:  ctx,
				info: &pb.GetCompileOutputRequest{PipelineUuid: pipelineId.String()},
			},
			want:    &pb.GetRunOutputResponse{Output: compileOutput},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			got, err := client.GetCompileOutput(tt.args.ctx, tt.args.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompileOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !strings.EqualFold(got.Output, tt.want.Output) {
					t.Errorf("GetCompileOutput() got = %v, want %v", got.Output, tt.want.Output)
				}
				if !reflect.DeepEqual(got.CompilationStatus, tt.want.CompilationStatus) {
					t.Errorf("GetCompileOutput() got = %v, want %v", got.CompilationStatus, tt.want.CompilationStatus)
				}
			}
		})
	}
}

func TestPlaygroundController_GetRunOutput(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewPlaygroundServiceClient(conn)
	pipelineMeta := pb.GetRunOutputRequest{
		PipelineUuid: uuid.NewString(),
	}
	runOutput, err := client.GetRunOutput(ctx, &pipelineMeta)
	if err != nil {
		t.Fatalf("runCode failed: %v", err)
	}
	log.Printf("Response: %+v", runOutput)
}
