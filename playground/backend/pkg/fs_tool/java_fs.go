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
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	parentBaseFileFolder = "pkg"

	javaBaseFileFolder     = parentBaseFileFolder + "/executable_files"
	javaSrcFileFolder      = javaBaseFileFolder + "/src"
	javaBinFileFolder      = javaBaseFileFolder + "/bin"
	javaExecutableFileType = "java"
	javaCompiledFileType   = "class"
)

type JavaFileSystemService struct{}

func (jFS *JavaFileSystemService) PrepareFolders() error {
	// Create pkg/executable_files/src folder
	err := os.MkdirAll(javaSrcFileFolder, fs.ModePerm)
	if err != nil {
		return err
	}

	// Create pkg/executable_files/bin folder
	err = os.MkdirAll(javaBinFileFolder, fs.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (jFS *JavaFileSystemService) DeletePreparedFolders() (bool, error) {
	// Check that pkg/executable_files exists
	if _, err := os.Stat(javaBaseFileFolder); os.IsNotExist(err) {
		return false, nil
	}
	// Delete pkg/executable_files folder and any children it contains (src/bin)
	err := os.RemoveAll(javaBaseFileFolder)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (jFS *JavaFileSystemService) CreateExecutableFile(code string, pipelineId int) (string, error) {
	// Check that pkg/executable_files/src folder exists
	srcFolderPath := javaSrcFileFolder
	if _, err := os.Stat(srcFolderPath); os.IsNotExist(err) {
		return "", err
	}

	fileName := getFileName(pipelineId, javaExecutableFileType)
	filePath := filepath.Join(srcFolderPath, fileName)
	err := ioutil.WriteFile(filePath, []byte(code), 0600)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (jFS *JavaFileSystemService) DeleteExecutableFile(pipelineId int) (bool, error) {
	fileName := getFileName(pipelineId, javaExecutableFileType)
	return deleteFile(fileName, javaSrcFileFolder)
}

func (jFS *JavaFileSystemService) DeleteCompiledFile(pipelineId int) (bool, error) {
	fileName := getFileName(pipelineId, javaCompiledFileType)
	return deleteFile(fileName, javaBinFileFolder)
}

// GetSrcPath returns path to the Java source folder.
func (jFS *JavaFileSystemService) GetSrcPath() string {
	return javaSrcFileFolder
}

// GetBinPath returns path to the Java bin folder.
func (jFS *JavaFileSystemService) GetBinPath() string {
	return javaBinFileFolder
}

// GetBaseFileFolderPath returns path to the base folder for Java files.
func (jFS *JavaFileSystemService) GetBaseFileFolderPath() string {
	return javaBaseFileFolder
}

// GetExecutableFileName returns path to the executable Java file by pipelineId.
func (jFS *JavaFileSystemService) GetExecutableFileName(pipelineId int) string {
	return getFileName(pipelineId, javaExecutableFileType)
}

// GetCompiledFileName returns path to the compiled Java file by pipelineId.
func (jFS *JavaFileSystemService) GetCompiledFileName(pipelineId int) string {
	return getFileName(pipelineId, javaCompiledFileType)
}

// deleteFile deletes file by fileName and folderPath.
func deleteFile(fileName string, folderPath string) (bool, error) {
	// Check that parent folder exists
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return false, nil
	}

	// Check that file exists
	filePath := filepath.Join(folderPath, fileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false, nil
	}

	err := os.Remove(filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

// getFileName returns fileName by pipelineId and fileType.
func getFileName(pipelineId int, fileType string) string {
	return fmt.Sprintf("%d.%s", pipelineId, fileType)
}
