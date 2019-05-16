package cmdflag

// Error defines the error type for this package.
type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNoCommand is returned when no command was found on the command line.
	ErrNoCommand Error = "no command specified"
	// ErrMissingCommandName is returned for an invalid command name (empty).
	ErrMissingCommandName Error = "missing command name"
	// ErrMissingInitializer is returned when a command does not have its initializer defined.
	ErrMissingInitializer Error = "missing command initializer"
	// ErrDuplicateCommand is returned when a command is redefined.
	ErrDuplicateCommand Error = "duplicated command"
)
