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

// Package executors
package executors

import (
	"beam.apache.org/playground/backend/internal/fs_tool"
	"bufio"
	"google.golang.org/grpc/grpclog"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	beamJarPath    = "/opt/apache/beam/jars/beam-sdks-java-harness.jar"
	runnerJarPath  = "/opt/apache/beam/jars/beam-runners-direct.jar"
	slf4jPath      = "/opt/apache/beam/jars/slf4j-jdk14.jar"
	javaExtension  = ".java"
	javaCompileCmd = "javac"
	javaRunCmd     = "java"
	binFolder      = "bin"
)

// NewJavaExecutor creates an executor with Go specifics
func NewJavaExecutor(fs *fs_tool.LifeCycle, javaValidators, javaPreparation *[]functionWithArgs) *Executor {
	compileArgs := []string{"-d", binFolder, "-classpath", beamJarPath}
	fullClassPath := strings.Join([]string{binFolder, beamJarPath, runnerJarPath, slf4jPath}, ":")
	runArgs := []string{"-cp", fullClassPath}
	if javaValidators == nil {
		v := make([]functionWithArgs, 0)
		javaValidators = &v
	}
	if javaPreparation == nil {
		v := make([]functionWithArgs, 0)
		javaPreparation = &v
	}
	path, _ := os.Getwd()

	exec := new(Executor)
	exec.validators = *javaValidators
	exec.preparation = *javaPreparation
	exec.relativeFilePath = fs.GetRelativeExecutableFilePath()
	exec.absoulteFilePath = fs.GetAbsoluteExecutableFilePath()
	exec.dirPath = filepath.Join(path, fs.Folder.BaseFolder)
	exec.compileName = javaCompileCmd
	exec.runName = javaRunCmd
	exec.compileArgs = compileArgs
	exec.runArgs = runArgs
	return exec
}

// GetJavaValidators return validators methods that needed for Java file
func GetJavaValidators() *[]functionWithArgs {
	validatorArgs := make([]interface{}, 1)
	validatorArgs[0] = javaExtension
	pathCheckerValidator := functionWithArgs{
		do:   fs_tool.CheckPathIsValid,
		args: validatorArgs,
	}
	validators := []functionWithArgs{pathCheckerValidator}
	return &validators
}

// GetJavaPreparation return validation methods that needed for Java file
func GetJavaPreparation() *[]functionWithArgs {
	publicClassModification := functionWithArgs{
		do: removePublicClassModifier,
	}
	additionalPackage := functionWithArgs{
		do: removeAdditionalPackage,
	}
	validators := []functionWithArgs{publicClassModification, additionalPackage}
	return &validators
}

func removePublicClassModifier(filePath string, args ...interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during open file: %s, err: %s", filePath, err.Error())
		return err
	}

	pathSlice := strings.Split(filePath, "/")
	fileName := pathSlice[len(pathSlice)-1]
	folderPath := ""
	pathSlice = pathSlice[:len(pathSlice)-1]
	for _, folder := range pathSlice {
		folderPath = filepath.Join(folderPath, folder)
	}
	tmp, err := os.Create(fileName)
	if err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during create new temporary file, err: %s", err.Error())
		return err
	}

	if err := transferWithReplace(file, tmp, "public class ", "class "); err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during replace and move data from original file to to temporary file, err: %s", err.Error())
		return err
	}

	if err := tmp.Close(); err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during Close temporary file, err: %s", err.Error())
		return err
	}
	if err := file.Close(); err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during Close original file, err: %s", err.Error())
		return err
	}

	if err := os.Rename(tmp.Name(), filePath); err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during rename temporary file, err: %s", err.Error())
		return err
	}
	return nil
}

func removeAdditionalPackage(filePath string, args ...interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		grpclog.Errorf("Preparation:  remove additional package: Error during open file: %s, err: %s", filePath, err.Error())
		return err
	}

	pathSlice := strings.Split(filePath, "/")
	fileName := pathSlice[len(pathSlice)-1]
	folderPath := ""
	pathSlice = pathSlice[:len(pathSlice)-1]
	for _, folder := range pathSlice {
		folderPath = filepath.Join(folderPath, folder)
	}
	tmp, err := os.Create(fileName)
	if err != nil {
		grpclog.Errorf("Preparation: remove additional package: Error during create new temporary file, err: %s", err.Error())
		return err
	}

	if err := transferWithReplace(file, tmp, `package ([\w]+\.)+[\w]+;`, ""); err != nil {
		grpclog.Errorf("Preparation: remove additional package: Error during move data from original file to to temporary file, err: %s", err.Error())
		return err
	}

	if err := tmp.Close(); err != nil {
		grpclog.Errorf("Preparation: remove additional package: Error during Close temporary file, err: %s", err.Error())
		return err
	}
	if err := file.Close(); err != nil {
		grpclog.Errorf("Preparation: remove additional package: Error during Close original file, err: %s", err.Error())
		return err
	}

	if err := os.Rename(tmp.Name(), filePath); err != nil {
		grpclog.Errorf("Preparation: remove additional package: Error during rename temporary file, err: %s", err.Error())
		return err
	}
	return nil
}

func transferWithReplace(reader io.Reader, writer io.Writer, pattern, new string) error {
	reg := regexp.MustCompile(pattern)
	scanner := bufio.NewScanner(reader)
	firstLine := true
	for scanner.Scan() {
		if !firstLine {
			if _, err := io.WriteString(writer, "\n"); err != nil {
				return err
			}
		}
		line := scanner.Text()
		matches := reg.FindAllString(line, -1)
		for _, str := range matches {
			line = strings.ReplaceAll(line, str, new)
		}
		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
		firstLine = false
	}
	return scanner.Err()
}
