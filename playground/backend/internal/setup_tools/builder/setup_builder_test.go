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

package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"

	pb "beam.apache.org/playground/backend/internal/api/v1"
	"beam.apache.org/playground/backend/internal/environment"
	"beam.apache.org/playground/backend/internal/executors"
	"beam.apache.org/playground/backend/internal/fs_tool"
)

const emptyFolder = "emptyFolder"

var pythonPaths *fs_tool.LifeCyclePaths
var pythonSdkEnv *environment.BeamEnvs
var javaLC *fs_tool.LifeCycle
var javaPaths *fs_tool.LifeCyclePaths
var javaSdkEnv *environment.BeamEnvs
var goPaths *fs_tool.LifeCyclePaths
var goSdkEnv *environment.BeamEnvs
var scioPaths *fs_tool.LifeCyclePaths
var scioSdkEnv *environment.BeamEnvs

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	os.Mkdir(emptyFolder, 0666)

	pipelineId := uuid.New()

	pythonLC, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_PYTHON, pipelineId, "")
	pythonPaths = &pythonLC.Paths

	javaLC, _ = fs_tool.NewLifeCycle(pb.Sdk_SDK_JAVA, pipelineId, "")
	javaPaths = &javaLC.Paths
	javaLC.CreateFolders()
	os.Create(filepath.Join(javaPaths.AbsoluteExecutableFilePath))
	os.Create(filepath.Join(javaPaths.AbsoluteSourceFilePath))

	goLC, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_GO, pipelineId, "")
	goPaths = &goLC.Paths

	scioLC, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_SCIO, pipelineId, "")
	scioPaths = &scioLC.Paths

	executorConfig := &environment.ExecutorConfig{
		CompileCmd:  "MOCK_COMPILE_CMD",
		CompileArgs: []string{"MOCK_COMPILE_ARG"},
		RunArgs:     []string{"MOCK_RUN_CMD"},
	}
	pythonSdkEnv = environment.NewBeamEnvs(pb.Sdk_SDK_PYTHON, executorConfig, "", 0)
	javaSdkEnv = environment.NewBeamEnvs(pb.Sdk_SDK_JAVA, executorConfig, "", 0)
	goSdkEnv = environment.NewBeamEnvs(pb.Sdk_SDK_GO, executorConfig, "", 0)
	scioSdkEnv = environment.NewBeamEnvs(pb.Sdk_SDK_SCIO, executorConfig, "", 0)
}

func teardown() {
	os.Remove(emptyFolder)
	javaLC.DeleteFolders()
}

func TestCompiler(t *testing.T) {
	wantJavaExecutor := executors.NewExecutorBuilder().
		WithCompiler().
		WithCommand(javaSdkEnv.ExecutorConfig.CompileCmd).
		WithWorkingDir(javaPaths.AbsoluteBaseFolderPath).
		WithArgs(javaSdkEnv.ExecutorConfig.CompileArgs).
		WithFileName(GetFirstFileFromFolder(javaPaths.AbsoluteSourceFileFolderPath))

	wantGoExecutor := executors.NewExecutorBuilder().
		WithCompiler().
		WithCommand(goSdkEnv.ExecutorConfig.CompileCmd).
		WithWorkingDir(goPaths.AbsoluteBaseFolderPath).
		WithArgs(goSdkEnv.ExecutorConfig.CompileArgs).
		WithFileName(goPaths.AbsoluteSourceFilePath)

	wantScioExecutor := executors.NewExecutorBuilder().
		WithCompiler().
		WithCommand(scioSdkEnv.ExecutorConfig.CompileCmd).
		WithWorkingDir(scioPaths.AbsoluteBaseFolderPath).
		WithArgs(scioSdkEnv.ExecutorConfig.CompileArgs).
		WithFileName(scioPaths.AbsoluteSourceFilePath)

	type args struct {
		paths  *fs_tool.LifeCyclePaths
		sdkEnv *environment.BeamEnvs
	}
	tests := []struct {
		name string
		args args
		want *executors.ExecutorBuilder
	}{
		{
			// Test case with calling Setup with correct data.
			// As a result, want to receive an expected compiler builder.
			name: "Test correct compiler builder with java sdk",
			args: args{
				paths:  javaPaths,
				sdkEnv: javaSdkEnv,
			},
			want: &wantJavaExecutor.ExecutorBuilder,
		},
		{
			name: "Test correct compiler builder with go sdk",
			args: args{
				paths:  goPaths,
				sdkEnv: goSdkEnv,
			},
			want: &wantGoExecutor.ExecutorBuilder,
		},
		{
			name: "Test correct compiler builder with scio sdk",
			args: args{
				paths:  scioPaths,
				sdkEnv: scioSdkEnv,
			},
			want: &wantScioExecutor.ExecutorBuilder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compiler(tt.args.paths, tt.args.sdkEnv)
			if !reflect.DeepEqual(fmt.Sprint(got.Build()), fmt.Sprint(tt.want.Build())) {
				t.Errorf("Compiler() = %v, want %v", got.Build(), tt.want.Build())
			}
		})
	}
}

func TestRunnerBuilder(t *testing.T) {
	incorrectPaths := *javaPaths
	incorrectPaths.AbsoluteExecutableFileFolderPath = emptyFolder
	incorrectPaths.AbsoluteBaseFolderPath = emptyFolder

	wantPythonExecutor := executors.NewExecutorBuilder().
		WithRunner().
		WithExecutableFileName(pythonPaths.AbsoluteExecutableFilePath).
		WithWorkingDir(pythonPaths.AbsoluteBaseFolderPath).
		WithCommand(pythonSdkEnv.ExecutorConfig.RunCmd).
		WithArgs(pythonSdkEnv.ExecutorConfig.RunArgs).
		WithPipelineOptions(strings.Split("", " "))

	arg := replaceLogPlaceholder(javaPaths, javaSdkEnv.ExecutorConfig)
	javaClassName, err := javaPaths.ExecutableName(javaPaths.AbsoluteExecutableFileFolderPath)
	if err != nil {
		panic(err)
	}
	wantJavaExecutor := executors.NewExecutorBuilder().
		WithRunner().
		WithExecutableFileName(javaClassName).
		WithWorkingDir(javaPaths.AbsoluteBaseFolderPath).
		WithCommand(javaSdkEnv.ExecutorConfig.RunCmd).
		WithArgs(arg).
		WithPipelineOptions(strings.Split("", " "))

	wantGoExecutor := executors.NewExecutorBuilder().
		WithRunner().
		WithWorkingDir(goPaths.AbsoluteBaseFolderPath).
		WithCommand(goPaths.AbsoluteExecutableFilePath).
		WithExecutableFileName("").
		WithArgs(goSdkEnv.ExecutorConfig.RunArgs).
		WithPipelineOptions(strings.Split("", " "))

	scioClassName, err := scioPaths.ExecutableName(scioPaths.AbsoluteBaseFolderPath)
	if err != nil {
		panic(err)
	}
	stringArg := fmt.Sprintf("%s %s %s", scioSdkEnv.ExecutorConfig.RunArgs[0], scioClassName, "")
	wantScioExecutor := executors.NewExecutorBuilder().
		WithRunner().
		WithWorkingDir(scioPaths.ProjectDir).
		WithCommand(scioSdkEnv.ExecutorConfig.RunCmd).
		WithArgs([]string{stringArg})

	type args struct {
		paths           *fs_tool.LifeCyclePaths
		pipelineOptions string
		sdkEnv          *environment.BeamEnvs
	}
	tests := []struct {
		name string
		args args
		want *executors.ExecutorBuilder
	}{
		{
			// Test case with calling Setup with correct data.
			// As a result, want to receive an expected run builder.
			name: "Test correct run builder with Python sdk",
			args: args{
				paths:  pythonPaths,
				sdkEnv: pythonSdkEnv,
			},
			want: &wantPythonExecutor.ExecutorBuilder,
		},
		{
			name: "Test correct run builder with Java sdk",
			args: args{
				paths:  javaPaths,
				sdkEnv: javaSdkEnv,
			},
			want: &wantJavaExecutor.ExecutorBuilder,
		},
		{
			name: "Test incorrect run builder with Java sdk",
			args: args{
				paths:  &incorrectPaths,
				sdkEnv: javaSdkEnv,
			},
			want: nil,
		},
		{
			name: "Test correct run builder with Go sdk",
			args: args{
				paths:  goPaths,
				sdkEnv: goSdkEnv,
			},
			want: &wantGoExecutor.ExecutorBuilder,
		},
		{
			name: "Test correct run builder with Scio sdk",
			args: args{
				paths:  scioPaths,
				sdkEnv: scioSdkEnv,
			},
			want: &wantScioExecutor.ExecutorBuilder,
		},
		{
			name: "Test incorrect run builder with Scio sdk",
			args: args{
				paths:  &incorrectPaths,
				sdkEnv: scioSdkEnv,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Runner(tt.args.paths, tt.args.pipelineOptions, tt.args.sdkEnv)
			if tt.want != nil {
				if !reflect.DeepEqual(fmt.Sprint(got.Build()), fmt.Sprint(tt.want.Build())) {
					t.Errorf("Runner() got = %v, want %v", got.Build(), tt.want.Build())
				}
			} else {
				if tt.want != got {
					t.Errorf("Runner() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTestRunner(t *testing.T) {
	incorrectJavaPaths := *javaPaths
	incorrectJavaPaths.AbsoluteExecutableFileFolderPath = emptyFolder

	className, err := javaPaths.ExecutableName(javaPaths.AbsoluteExecutableFileFolderPath)
	if err != nil {
		panic(err)
	}
	wantJavaExecutor := executors.NewExecutorBuilder().
		WithTestRunner().
		WithExecutableFileName(className).
		WithCommand(javaSdkEnv.ExecutorConfig.TestCmd).
		WithArgs(javaSdkEnv.ExecutorConfig.TestArgs).
		WithWorkingDir(javaPaths.AbsoluteBaseFolderPath)

	wantGoExecutor := executors.NewExecutorBuilder().
		WithTestRunner().
		WithExecutableFileName(goPaths.AbsoluteSourceFileFolderPath).
		WithCommand(javaSdkEnv.ExecutorConfig.TestCmd).
		WithArgs(javaSdkEnv.ExecutorConfig.TestArgs).
		WithWorkingDir(goPaths.AbsoluteSourceFileFolderPath)

	wantPythonExecutor := executors.NewExecutorBuilder().
		WithTestRunner().
		WithExecutableFileName(pythonPaths.AbsoluteExecutableFilePath).
		WithCommand(pythonSdkEnv.ExecutorConfig.TestCmd).
		WithArgs(pythonSdkEnv.ExecutorConfig.TestArgs).
		WithWorkingDir(pythonPaths.AbsoluteSourceFileFolderPath)

	type args struct {
		paths  *fs_tool.LifeCyclePaths
		sdkEnv *environment.BeamEnvs
	}
	tests := []struct {
		name string
		args args
		want *executors.ExecutorBuilder
	}{
		{
			name: "Test correct run builder with Java sdk",
			args: args{
				paths:  javaPaths,
				sdkEnv: javaSdkEnv,
			},
			want: &wantJavaExecutor.ExecutorBuilder,
		},
		{
			name: "Test incorrect run builder with Java sdk",
			args: args{
				paths:  &incorrectJavaPaths,
				sdkEnv: javaSdkEnv,
			},
			want: nil,
		},
		{
			name: "Test correct run builder with GO sdk",
			args: args{
				paths:  goPaths,
				sdkEnv: goSdkEnv,
			},
			want: &wantGoExecutor.ExecutorBuilder,
		},
		{
			name: "Test correct run builder with Python sdk",
			args: args{
				paths:  pythonPaths,
				sdkEnv: pythonSdkEnv,
			},
			want: &wantPythonExecutor.ExecutorBuilder,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := TestRunner(tt.args.paths, tt.args.sdkEnv)
			if tt.want != nil {
				if !reflect.DeepEqual(fmt.Sprint(got.Build()), fmt.Sprint(tt.want.Build())) {
					t.Errorf("TestRunner() got = %v, want %v", got.Build(), tt.want.Build())
				}
			} else {
				if tt.want != got {
					t.Errorf("TestRunner() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_replaceLogPlaceholder(t *testing.T) {
	execConfig := *pythonSdkEnv.ExecutorConfig
	execConfig.RunArgs = []string{"arg1", javaLogConfigFilePlaceholder}
	type args struct {
		paths          *fs_tool.LifeCyclePaths
		executorConfig *environment.ExecutorConfig
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test to check the replacement of log work with prepared running args",
			args: args{
				paths:          pythonPaths,
				executorConfig: &execConfig,
			},
			want: []string{"arg1", filepath.Join(pythonPaths.AbsoluteBaseFolderPath, javaLogConfigFileName)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceLogPlaceholder(tt.args.paths, tt.args.executorConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replaceLogPlaceholder() = %v, want %v", got, tt.want)
			}
		})
	}
}
