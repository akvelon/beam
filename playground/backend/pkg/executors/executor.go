package executors

type executor interface {
	//Validate executable file
	Validate(filePath string) (bool, error)

	//Run executable file and write logs to writer
	Run(filePath string) (string, error)
}
