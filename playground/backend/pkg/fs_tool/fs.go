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

package fs_tool

import (
	pb "beam.apache.org/playground/backend/pkg/api"
	"errors"
	"fmt"
)

// FileSystemService interface for all type(Java/Go/Python/Scio) of FileSystem services.
type FileSystemService interface {

	// PrepareFolders creates folders to save a file with the code for the selected sdk in them.
	PrepareFolders() error

	// DeletePreparedFolders deletes folders which were created for executable files.
	// Returns true if parent folders existed and were deleted
	// Returns false if parent folder doesn't exist
	DeletePreparedFolders() (bool, error)

	// CreateExecutableFile creates the file and puts into it the given code from parameters.
	CreateExecutableFile(code string, pipelineId int) (string, error)

	// DeleteExecutableFile deletes the executable file.
	// Returns true if the executable file existed and was deleted
	// Returns false if executable file doesn't exist
	DeleteExecutableFile(pipelineId int) (bool, error)

	// DeleteCompiledFile deletes the file, which was generated after the compilation step.
	// Returns true if the compiled file existed and was deleted
	// Returns false if the compiled file doesn't exist
	DeleteCompiledFile(pipelineId int) (bool, error)
}

// NewFileSystemService returns a corresponding FileSystemService depending on the given sdk.
func NewFileSystemService(sdk pb.Sdk) (FileSystemService, error) {
	switch sdk {
	case pb.Sdk_JAVA:
		return &JavaFileSystemService{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("%s isn't supported now", sdk))
	}
}
