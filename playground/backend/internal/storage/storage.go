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

type storage interface {
	// SetCompileOutput set output of compilation to a storage
	SetCompileOutput(pipelineUuid string, output string) error

	// GetCompileOutput get output of compilation by pipelineUuid
	GetCompileOutput(pipelineUuid string) (string, error)

	// SetRunOutput set output of execution to a storage
	SetRunOutput(pipelineUuid string, output string) error

	// GetRunOutput get output of execution by pipelineUuid
	GetRunOutput(pipelineUuid string) (string, error)

	// SetStatusPipeline set pipeline status to a storage
	SetStatusPipeline(pipelineUuid string, status *pb.CheckStatusResponse) error

	// GetStatusPipeline get status of pipeline by pipelineUuid
	GetStatusPipeline(pipelineUuid string) *pb.CheckStatusResponse
}
