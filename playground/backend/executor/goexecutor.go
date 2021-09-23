package executor

type GoExecutor struct{}

func (goExec GoExecutor) Validate(filePath string) (bool, error) {
	return true, nil
}

func (goExec GoExecutor) Run(filePath string) (string, error) {
	return "", nil
}
