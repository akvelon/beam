package executor

type JavaExecutor struct{}

func (javaExec JavaExecutor) Validate(filePath string) (bool, error) {
	return true, nil
}

func (javaExec JavaExecutor) Run(filePath string) (string, error) {
	return "", nil
}
