package preparators

import (
	playground "beam.apache.org/playground/backend/internal/api"
	"beam.apache.org/playground/backend/internal/fs_tool"
	"github.com/google/uuid"
	"os"
	"strings"
	"testing"
)

func Test_removePublicClassModifier(t *testing.T) {
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
			args:    args{[]interface{}{"someFile.java"}},
			wantErr: true,
		},
		{
			name:     "original file exists",
			args:     args{[]interface{}{lc.GetAbsoluteExecutableFilePath()}},
			wantCode: codeWithoutPublicClass,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removePublicClassModifier(tt.args.args...); (err != nil) != tt.wantErr {
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

func Test_removeAdditionalPackage(t *testing.T) {
	codeWithPackage := "package org.some.package;\n\nclass Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"
	codeWithoutPackage := "\n\nclass Class {\n    public static void main(String[] args) {\n        System.out.println(\"Hello World!\");\n    }\n}"

	lc, _ := fs_tool.NewLifeCycle(playground.Sdk_SDK_JAVA, uuid.New())
	_ = lc.CreateFolders()
	defer lc.DeleteFolders()
	_, _ = lc.CreateExecutableFile(codeWithPackage)

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
			args:    args{[]interface{}{"someFile.java"}},
			wantErr: true,
		},
		{
			name:     "original file exists",
			args:     args{[]interface{}{lc.GetAbsoluteExecutableFilePath()}},
			wantCode: codeWithoutPackage,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := removeAdditionalPackage(tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("removeAdditionalPackage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				data, err := os.ReadFile(tt.args.args[0].(string))
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
