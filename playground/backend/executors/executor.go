package executors

import (
	"io"
)

type executor interface {
	//Validate executable file
	Validate(filePath string) bool

	//Run executable file and write logs to writer
	Run(filePath string, writer io.Writer) error
}

type GoExecutor struct{}

func (goExec GoExecutor) Validate(filePath string) bool {
	return true
}

func (goExec GoExecutor) Run(filePath string, writer io.Writer) error {
	return nil
}

type JavaExecutor struct{}

func (javaExec JavaExecutor) Validate(filePath string) bool {
	return true
}

func (javaExec JavaExecutor) Run(filePath string, writer io.Writer) error {
	return nil
}
