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
func NewJavaExecutor(fs *fs_tool.LifeCycle, javaValidators *[]validatorWithArgs, javaPreparation *[]preparationWithArgs) *Executor {
	compileArgs := []string{"-d", binFolder, "-classpath", beamJarPath}
	fullClassPath := strings.Join([]string{binFolder, beamJarPath, runnerJarPath, slf4jPath}, ":")
	runArgs := []string{"-cp", fullClassPath}
	if javaValidators == nil {
		v := make([]validatorWithArgs, 0)
		javaValidators = &v
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
func GetJavaValidators() *[]validatorWithArgs {
	validatorArgs := make([]interface{}, 1)
	validatorArgs[0] = javaExtension
	pathCheckerValidator := validatorWithArgs{
		validator: fs_tool.CheckPathIsValid,
		args:      validatorArgs,
	}
	validators := []validatorWithArgs{pathCheckerValidator}
	return &validators
}

// GetJavaPreparation return validation methods that needed for Java file
func GetJavaPreparation() *[]preparationWithArgs {
	publicClassModification := preparationWithArgs{
		prepare: removePublicClassModification,
	}
	additionalPackage := preparationWithArgs{
		prepare: removeAdditionalPackage,
	}
	validators := []preparationWithArgs{publicClassModification, additionalPackage}
	return &validators
}

func removeAdditionalPackage(filePath string, args ...interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		grpclog.Errorf("Preparation:  remove additional package: Error during open file: %s", filePath)
		return err
	}

	tmp, err := os.Create("temporary.java")
	if err != nil {
		grpclog.Error("Preparation: remove additional package: Error during create new temporary file")
		return err
	}

	if err := removePackageString(file, tmp); err != nil {
		grpclog.Error("Preparation: remove additional package: Error during move data from original file to to temporary file")
		return err
	}

	if err := tmp.Close(); err != nil {
		grpclog.Error("Preparation: remove additional package: Error during Close temporary file")
		return err
	}
	if err := file.Close(); err != nil {
		grpclog.Error("Preparation: remove additional package: Error during Close original file")
		return err
	}

	pathSlice := strings.Split(filePath, "/")
	fileName := pathSlice[len(pathSlice)-1]
	if err := os.Rename(tmp.Name(), fileName); err != nil {
		grpclog.Error("Preparation: remove additional package: Error during rename temporary file")
		return err
	}
	return nil
}

func removePackageString(reader io.Reader, writer io.Writer) error {
	reg := regexp.MustCompile(`package ([\w]+\.)+[\w]+;`)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		matches := reg.FindAllString(line, -1)
		for _, str := range matches {
			line = strings.ReplaceAll(line, str, "")
		}
		if _, err := io.WriteString(writer, line+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}

func removePublicClassModification(filePath string, args ...interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		grpclog.Errorf("Preparation: remove public modification for class: Error during open file: %s", filePath)
		return err
	}

	tmp, err := os.Create("temporary.java")
	if err != nil {
		grpclog.Error("Preparation: remove public modification for class: Error during create new temporary file")
		return err
	}

	if err := replace(file, tmp); err != nil {
		grpclog.Error("Preparation: remove public modification for class: Error during replace and move data from original file to to temporary file")
		return err
	}

	if err := tmp.Close(); err != nil {
		grpclog.Error("Preparation: remove public modification for class: Error during Close temporary file")
		return err
	}
	if err := file.Close(); err != nil {
		grpclog.Error("Preparation: remove public modification for class: Error during Close original file")
		return err
	}

	pathSlice := strings.Split(filePath, "/")
	fileName := pathSlice[len(pathSlice)-1]
	if err := os.Rename(tmp.Name(), fileName); err != nil {
		grpclog.Error("Preparation: remove public modification for class: Error during rename temporary file")
		return err
	}
	return nil
}

func replace(reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "public class ", "class ")
		if _, err := io.WriteString(writer, line+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
