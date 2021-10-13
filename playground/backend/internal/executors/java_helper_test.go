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

package executors

import (
	pb "beam.apache.org/playground/backend/internal/api"
	"beam.apache.org/playground/backend/internal/fs_tool"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
)

var (
	javaExecutor *Executor
	pipelineId   = uuid.New()
)

const (
	javaCode = "class HelloWorld {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"
)

func TestMain(m *testing.M) {
	javaFS := setup()
	defer teardown(javaFS)
	m.Run()
}

func setup() *fs_tool.LifeCycle {
	javaFS, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_JAVA, pipelineId)
	_ = javaFS.CreateFolders()
	_, _ = javaFS.CreateExecutableFile(javaCode)
	javaExecutor = NewJavaExecutor(javaFS, GetJavaValidators(), GetJavaPreparation())
	return javaFS
}

func teardown(javaFS *fs_tool.LifeCycle) {
	err := javaFS.DeleteFolders()
	if err != nil {
		return
	}
}

func TestValidateJavaFile(t *testing.T) {
	err := javaExecutor.Validate()
	if err != nil {
		t.Fatalf(`TestValidateJavaFile error: %v `, err)
	}
}

func TestCompileJavaFile(t *testing.T) {
	err := javaExecutor.Compile()
	if err != nil {
		t.Fatalf("TestCompileJavaFile: Unexpexted error at compiliation: %s ", err.Error())
	}
}

func TestRunJavaFile(t *testing.T) {
	className := "HelloWorld"
	expected := "Hello World!\n"
	out, err := javaExecutor.Run(className)
	if expected != out || err != nil {
		t.Fatalf(`TestRunJavaFile: '%q, %v' doesn't match for '%#q', nil`, out, err, expected)
	}
}

func Test_removeAdditionalPackage(t *testing.T) {
	codeWithPackage := "package org.some.package;\n\nclass Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"
	codeWithoutPackage := "\n\nclass Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"

	javaFS, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_JAVA, pipelineId)
	_, err := javaFS.CreateExecutableFile(codeWithPackage)
	if err != nil {
		t.Errorf("removeAdditionalPackage() error during preparing files for test, err: %s", err)
		return
	}

	type args struct {
		filePath string
		args     []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantCode string
		wantErr  bool
	}{
		{
			name: "original file doesn't exist",
			args: args{
				filePath: "someFile.java",
				args:     make([]interface{}, 0),
			},
			wantErr: true,
		},
		{
			name: "original file exists",
			args: args{
				filePath: javaFS.GetAbsoluteExecutableFilePath(),
				args:     make([]interface{}, 0),
			},
			wantCode: codeWithoutPackage,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeAdditionalPackage(tt.args.filePath, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("removeAdditionalPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				data, err := os.ReadFile(tt.args.filePath)
				if err != nil {
					t.Errorf("removeAdditionalPackage() unexpected error = %v", err)
				}
				if !strings.EqualFold(string(data), tt.wantCode) {
					t.Errorf("removeAdditionalPackage() code = {%v}, wantCode {%v}", string(data), tt.wantCode)
				}
			}
		})
	}
}

func Test_removePublicClassModifier(t *testing.T) {
	codeWithPublicClass := "public class Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"
	codeWithoutPublicClass := "class Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"

	javaFS, _ := fs_tool.NewLifeCycle(pb.Sdk_SDK_JAVA, pipelineId)
	_, err := javaFS.CreateExecutableFile(codeWithPublicClass)
	if err != nil {
		t.Errorf("removePublicClassModifie() error during preparing files for test, err: %s", err)
		return
	}

	type args struct {
		filePath string
		args     []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantCode string
		wantErr  bool
	}{
		{
			name: "original file doesn't exist",
			args: args{
				filePath: "someFile.java",
				args:     make([]interface{}, 0),
			},
			wantErr: true,
		},
		{
			name: "original file exists",
			args: args{
				filePath: javaFS.GetAbsoluteExecutableFilePath(),
				args:     make([]interface{}, 0),
			},
			wantCode: codeWithoutPublicClass,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removePublicClassModifier(tt.args.filePath, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("removePublicClassModifier() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				data, err := os.ReadFile(tt.args.filePath)
				if err != nil {
					t.Errorf("removePublicClassModifier() unexpected error = %v", err)
				}
				if !strings.EqualFold(string(data), tt.wantCode) {
					t.Errorf("removePublicClassModifier() code = {%v}, wantCode {%v}", string(data), tt.wantCode)
				}
			}
		})
	}
}
