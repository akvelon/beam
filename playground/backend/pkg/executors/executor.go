package executors

// executor Interface for all executors (Java/Python/Go/SCIO)
type executor interface {
	// Validate validates executable file.
	// Return result of validation (true/false) and error if it occurs
	Validate(filePath string) (bool, error)

	// Compile compiles executable file.
	// Return error if it occurs
	Compile(filePath string) error

	// Run runs executable file.
	// Return logs and error if it occurs
	Run(filePath string) (string, error)
}
