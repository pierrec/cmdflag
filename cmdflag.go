// Package cmdflag provides simple command line commands processing
// on top of the standard library flag package.
//
// It strives to be lightweight (only relying on the standard library) and
// fits naturally with the usage of the flag package.
package cmdflag

import (
	"flag"
	"io"
	"os"
	"sync"
)

type (
	// Application defines the attributes of a Command.
	Application struct {
		Name  string                      // Command name
		Descr string                      // Short description
		Args  string                      // Description of the expected arguments
		Help  string                      // Displayed when used with the help command
		Err   flag.ErrorHandling          // Arguments error handling
		Init  func(*flag.FlagSet) Handler // Initialize the arguments when the command is matched
	}

	// Handler is the function called when a matching command is found.
	// It returns the number of arguments consumed or an error.
	Handler func(args ...string) (int, error)

	// Command represents a command line command.
	Command struct {
		fset *flag.FlagSet
		mu   sync.Mutex
		subs []*Command // Commands supported by this command

		Application
		// Usage is the function used to display the usage description.
		Usage func()
	}
)

// New instantiates the top level command based on the provided flag set.
// If no flag set is supplied, then it defaults to flag.CommandLine.
func New(fset *flag.FlagSet) *Command {
	if fset == nil {
		fset = flag.CommandLine
		fset.Usage = nil // Remove the default usage as it does not know about commands
	}
	return &Command{
		Application: Application{Name: program()},
		fset:        fset,
	}
}

// AddHelp adds a help command to display additional information for commands.
func (c *Command) AddHelp() error {
	return addHelpCommand(c)
}

// MustAddHelp is similar to AddHelp but panics if an error is encountered.
func (c *Command) MustAddHelp() {
	err := addHelpCommand(c)
	if err != nil {
		panic(err)
	}
}

// Add adds a new command with its name and description and returns the new command.
//
// It is safe to be called from multiple go routines (typically in init functions).
//
// The command initializer is called only when the command is present on the command line.
// The handler is called with the remaining arguments once the command flags have been parsed successfully.
//
// Command names must be unique and non empty.
func (c *Command) Add(app Application) (*Command, error) {
	if app.Name == "" {
		return nil, ErrMissingCommandName
	}
	if app.Init == nil {
		return nil, ErrMissingInitializer
	}
	sub := &Command{Application: app}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, sub := range c.subs {
		if name := sub.Application.Name; name == app.Name {
			return nil, ErrDuplicateCommand
		}
	}
	c.subs = append(c.subs, sub)
	return sub, nil
}

// MustAdd is similar to Add but panics if an error is encountered.
func (c *Command) MustAdd(app Application) *Command {
	cc, err := c.Add(app)
	if err != nil {
		panic(err)
	}
	return cc
}

// Commands returns all the commands defined on c.
func (c *Command) Commands() []*Command {
	return c.subs
}

// Output returns the output used for usage. It defaults to os.Stderr.
func (c *Command) Output() io.Writer {
	return fsetOutput(c.fset)
}

// Parse parses the command line arguments from the argument list, which should not include the command name
// and including the global flags and, if any, the command and its flags.
//
// To be called in lieu of flag.Parse().
//
// If no arguments are supplied, it defaults to os.Args[1:].
// If the VersionBoolFlag is defined as a global boolean flag, then the program version is displayed and the program
// stops.
// If the FullVersionBoolFlag is defined as a global boolean flag, then the full program version is displayed and
// the program stops.
func (c *Command) Parse(args ...string) error {
	if args == nil {
		args = os.Args[1:]
	}
	fset := c.fset
	out := fsetOutput(fset)

	// Usage: use the one supplied to c or its fset, or the default one.
	if c.Usage != nil {
		fset.Usage = c.Usage
	} else if fset.Usage == nil {
		fset.Usage = usage(out, c)
	}

	// Global flags.
	if err := fset.Parse(args); err != nil {
		return err
	}

	// Handle builtin flags.
	if hasBoolFlag(fset, VersionBoolFlag) {
		version(out)
		return nil
	}
	if hasBoolFlag(fset, FullVersionBoolFlag) {
		fullversion(out)
		return nil
	}

	// Only error on the first level.
	return c.run(0, fset, true)
}

// run a command and its own ones recursively.
func (c *Command) run(start int, fset *flag.FlagSet, doerror bool) error {
	// No command.
	if fset.NArg() == 0 || len(c.subs) == 0 {
		return nil
	}

	out := fsetOutput(fset)
	args := fset.Args()[start:]
	s := args[0]
	args = args[1:]
	for _, sub := range c.subs {
		if sub.Application.Name != s {
			continue
		}

		fs := flag.NewFlagSet("", sub.Application.Err)
		fs.SetOutput(out)
		fs.Usage = usage(out, sub)
		sub.fset = fs
		handler := sub.Application.Init(fs)
		// Command specific arguments.
		if err := fs.Parse(args); err != nil {
			return err
		}
		// Command handler.
		n, err := handler(args[len(args)-fs.NArg():]...)
		if err != nil {
			return err
		}
		// Next command.
		return sub.run(n, fs, false)
	}

	if doerror {
		return ErrNoCommand
	}

	return nil
}
