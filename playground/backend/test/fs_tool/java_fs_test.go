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
	"beam.apache.org/playground/backend/pkg/fs_tool"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"
)

const (
	pipelineId     = 1000
	baseFolderPath = "pkg/executable_files"
	binFolderPath  = baseFolderPath + "/bin"
	srcFolderPath  = baseFolderPath + "/src"
)

func TestPrepareFolders(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	err := javaFileSystemService.PrepareFolders()

	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if _, err = os.Stat(binFolderPath); os.IsNotExist(err) {
		t.Errorf("%s folder was not created", binFolderPath)
	}
	if _, err = os.Stat(srcFolderPath); os.IsNotExist(err) {
		t.Errorf("%s folder was not created", srcFolderPath)
	}
	_ = os.RemoveAll("pkg")
}

func TestDeletePreparedFoldersWhenNoOneFolderExists(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeletePreparedFolders()
	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if flag {
		t.Errorf("method should return false because folders were not existed")
	}
}

func TestDeletePreparedFolders(t *testing.T) {
	// Create src folder
	createSrcFolder(t)
	// Create bin folder
	createBinFolder(t)
	defer os.RemoveAll("pkg")

	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeletePreparedFolders()

	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if !flag {
		t.Errorf("method should return true because folders were existed")
	}
	// check that all folders were deleted
	if _, err = os.Stat(srcFolderPath); !os.IsNotExist(err) {
		t.Errorf("%s folder should be deleted", srcFolderPath)
	}
	if _, err = os.Stat(binFolderPath); !os.IsNotExist(err) {
		t.Errorf("%s folder should be deleted", binFolderPath)
	}
	if _, err = os.Stat(baseFolderPath); !os.IsNotExist(err) {
		t.Errorf("%s folder should be deleted", baseFolderPath)
	}
}

func TestCreateExecutableFileWhenSrcFolderIsNotExist(t *testing.T) {
	code := "testCodeString"
	expectedErrMessage := fmt.Sprintf("stat %s: no such file or directory", srcFolderPath)
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	_, err := javaFileSystemService.CreateExecutableFile(code, pipelineId)
	if err == nil {
		t.Errorf("method should return error because of the %s folder is not exist", srcFolderPath)
	}
	if !strings.EqualFold(err.Error(), expectedErrMessage) {
		t.Errorf("unexpected error. expected: %s, actual: %s", expectedErrMessage, err.Error())
	}
}

func TestCreateExecutableFile(t *testing.T) {
	// Create src folder
	createSrcFolder(t)
	defer os.RemoveAll("pkg")

	code := "testCodeString"
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	filePath, err := javaFileSystemService.CreateExecutableFile(code, pipelineId)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("error during reading data from file %s: %s", filePath, err)
	}
	if !strings.EqualFold(code, string(data)) {
		t.Errorf("unexpected data into generated file. expected: %s, actual: %s", code, string(data))
	}
}

func TestDeleteExecutableFileWhenSrcFolderIsNotExist(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteExecutableFile(pipelineId)

	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if flag {
		t.Errorf("method should return false because folder was not existed")
	}
}

func TestDeleteExecutableFileWhenExecutableFileIsNotExist(t *testing.T) {
	createSrcFolder(t)
	defer os.RemoveAll("pkg")

	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteExecutableFile(pipelineId)

	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if flag {
		t.Errorf("method should return false because file was not existed")
	}
}

func TestDeleteExecutableFile(t *testing.T) {
	filePath := fmt.Sprintf("%s/%d.java", srcFolderPath, pipelineId)
	createSrcFolder(t)
	_, err := os.Create(filePath)
	if err != nil {
		t.Errorf("error during creating file %s: %s", filePath, err)
	}
	defer os.RemoveAll("pkg/")

	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteExecutableFile(pipelineId)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if !flag {
		t.Errorf("method should return true because file was existed")
	}
	if _, err = os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("%s file should be deleted", filePath)
	}
}

func TestDeleteCompiledFileWhenBinFolderIsNotExist(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteCompiledFile(pipelineId)

	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if flag {
		t.Errorf("method should return false because folder was not existed")
	}
}

func TestDeleteCompiledFileWhenCompiledFileIsNotExist(t *testing.T) {
	createBinFolder(t)
	defer os.RemoveAll("pkg")

	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteCompiledFile(pipelineId)

	if err != nil {
		t.Errorf("unexpecrted error: %s", err)
	}
	if flag {
		t.Errorf("method should return false because file was not existed")
	}
}

func TestDeleteCompiledFile(t *testing.T) {
	filePath := fmt.Sprintf("%s/%d.class", binFolderPath, pipelineId)
	createBinFolder(t)
	_, err := os.Create(filePath)
	if err != nil {
		t.Errorf("error during creating file %s: %s", filePath, err)
	}
	defer os.RemoveAll("pkg")

	javaFileSystemService := fs_tool.JavaFileSystemService{}

	flag, err := javaFileSystemService.DeleteCompiledFile(pipelineId)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if !flag {
		t.Errorf("method should return true because file was existed")
	}
	if _, err = os.Stat(filePath); !os.IsNotExist(err) {
		t.Errorf("%s file should be deleted", filePath)
	}
}

func TestGetJavaSrcPath(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}
	expectedPath := srcFolderPath
	actualPath := javaFileSystemService.GetSrcPath()
	if !strings.EqualFold(expectedPath, actualPath) {
		t.Errorf("unexpected result. expected: %s, actual: %s", expectedPath, actualPath)
	}
}

func TestGetJavaBinPath(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}
	expectedPath := binFolderPath
	actualPath := javaFileSystemService.GetBinPath()
	if !strings.EqualFold(expectedPath, actualPath) {
		t.Errorf("unexpected result. expected: %s, actual: %s", expectedPath, actualPath)
	}
}

func TestGetJavaBaseFileFolderPath(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}
	expectedPath := baseFolderPath
	actualPath := javaFileSystemService.GetBaseFileFolderPath()
	if !strings.EqualFold(expectedPath, actualPath) {
		t.Errorf("unexpected result. expected: %s, actual: %s", expectedPath, actualPath)
	}
}

func TestGetJavaExecutableFilePath(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}
	expectedPath := fmt.Sprintf("%d.java", pipelineId)
	actualPath := javaFileSystemService.GetExecutableFileName(pipelineId)
	if !strings.EqualFold(expectedPath, actualPath) {
		t.Errorf("unexpected result. expected: %s, actual: %s", expectedPath, actualPath)
	}
}

func TestGetJavaCompiledFilePath(t *testing.T) {
	javaFileSystemService := fs_tool.JavaFileSystemService{}
	expectedPath := fmt.Sprintf("%d.class", pipelineId)
	actualPath := javaFileSystemService.GetCompiledFileName(pipelineId)
	if !strings.EqualFold(expectedPath, actualPath) {
		t.Errorf("unexpected result. expected: %s, actual: %s", expectedPath, actualPath)
	}
}

// createSrcFolder create srcFolderPath folder.
func createSrcFolder(t *testing.T) {
	err := os.MkdirAll(srcFolderPath, fs.ModePerm)
	if err != nil {
		t.Errorf("couldn't create folder %s for test", srcFolderPath)
	}
}

// createBinFolder create binFolderPath folder.
func createBinFolder(t *testing.T) {
	err := os.MkdirAll(binFolderPath, fs.ModePerm)
	if err != nil {
		t.Errorf("couldn't create folder %s for test", binFolderPath)
	}
}
