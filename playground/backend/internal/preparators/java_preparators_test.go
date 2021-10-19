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

package preparators

import (
	playground "beam.apache.org/playground/backend/internal/api"
	"beam.apache.org/playground/backend/internal/fs_tool"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
)

func Test_replace(t *testing.T) {
	codeWithPublicClass := "public class Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"
	codeWithoutPublicClass := "class Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"

	lc, _ := fs_tool.NewLifeCycle(playground.Sdk_SDK_JAVA, uuid.New())
	_ = lc.CreateFolders()
	defer lc.DeleteFolders()
	_, _ = lc.CreateExecutableFile(codeWithPublicClass)

	type args struct {
		args []interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantCode string
		wantErr  bool
	}{
		{
			name:    "original file doesn't exist",
			args:    args{[]interface{}{"someFile.java", classWithPublicModifierPattern, classWithoutPublicModifierPattern}},
			wantErr: true,
		},
		{
			name:     "original file exists",
			args:     args{[]interface{}{lc.GetAbsoluteExecutableFilePath(), classWithPublicModifierPattern, classWithoutPublicModifierPattern}},
			wantCode: codeWithoutPublicClass,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := replace(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("removePublicClassModifier() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				data, err := os.ReadFile(tt.args.args[0].(string))
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

func TestGetJavaPreparation(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "all success",
			args: args{"MOCK_FILEPATH"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetJavaPreparation(tt.args.filePath); len(*got) != tt.want {
				t.Errorf("GetJavaPreparation() returns %v Preparators, want %v", len(*got), tt.want)
			}
		})
	}
}
