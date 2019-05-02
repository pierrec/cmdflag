// Package cmdflag provides simple command line commands processing
// on top of the standard library flag package.
//
// It strives to be lightweight (only relying on the standard library) and
// fit naturally with the usage of the flag package.
package cmdflag

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
)

// Usage is the function used to display the help message.
var Usage = func() {
	fset := flag.CommandLine
	out := fset.Output()

	program := strings.TrimSuffix(os.Args[0], ".exe")
	_, _ = fmt.Fprintf(out, "Usage of %s:\n", program)
	fset.PrintDefaults()

	_, _ = fmt.Fprintf(out, "\nSubcommands:")
	for _, c := range CommandLine.Commands() {
		app := c.Application
		_, _ = fmt.Fprintf(out, "\n%s\n%s %s\n", app.Descr, app.Name, app.Args)
		fs := flag.NewFlagSet(app.Name, app.Err)
		fs.PrintDefaults()
	}
}

// CommandLine is the top level command.
var CommandLine Command

type (
	// Application defines the attributes of a Command.
	Application struct {
		Name  string                          // Command name
		Descr string                          // Short description
		Args  string                          // Description of the expected arguments
		Help  string                          // Displayed when used with the help command
		Err   flag.ErrorHandling              // Arguments error handling
		Init  func(set *flag.FlagSet) Handler // Initialize the arguments when the command is matched
	}

	// Handler is the function called when a matching command is found.
	Handler func(...string) error

	// Command represents a command line command.
	Command struct {
		mu   sync.Mutex
		subs []*Command

		Application
	}
)

// AddHelp adds a help command to display additional information for commands.
func (c *Command) AddHelp() error {
	return addHelpCommand(c)
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
		return nil, fmt.Errorf("missing command name")
	}
	if app.Init == nil {
		return nil, fmt.Errorf("missing command initializer")
	}
	sub := &Command{Application: app}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, sub := range c.subs {
		if name := sub.Application.Name; name == app.Name {
			return nil, fmt.Errorf("command %s redeclared", name)
		}
	}
	c.subs = append(c.subs, sub)
	return sub, nil
}

// Commands returns the commands defined on c.
func (c *Command) Commands() []*Command {
	return c.subs
}

// Parse parses the command line arguments including the global flags and, if any, the command and its flags.
//
// To be called in lieu of flag.Parse().
//
// If the VersionBoolFlag is defined as a global boolean flag, then the program version is displayed and the program stops.
func Parse() error {
	args := os.Args
	if len(args) == 1 {
		return nil
	}

	// Global flags.
	fset := flag.CommandLine
	fset.Usage = Usage
	out := fsetOutput(fset)

	if err := fset.Parse(args[1:]); err != nil {
		return err
	}

	// Handle builtin flags.
	if hasBoolFlag(fset, VersionBoolFlag) {
		program := strings.TrimSuffix(args[0], ".exe")
		_, _ = fmt.Fprintf(out, "%s version %s %s/%s\n",
			program, buildinfo(),
			runtime.GOOS, runtime.GOARCH)
		return nil
	}
	if hasBoolFlag(fset, FullVersionBoolFlag) {
		program := strings.TrimSuffix(args[0], ".exe")
		_, _ = fmt.Fprintf(out, "%s full version %s %s/%s compiled by %s (%s)\n%s\n",
			program, buildinfo(),
			runtime.GOOS, runtime.GOARCH,
			runtime.Compiler, runtime.Version(),
			fullbuildinfo())
		return nil
	}

	// Only error on the first level.
	return run(&CommandLine, args, fset, true)
}

// run a command and its own ones recursively.
func run(c *Command, args []string, fset *flag.FlagSet, doerror bool) error {
	// No command.
	if fset.NArg() == 0 {
		return nil
	}

	out := fsetOutput(fset)
	idx := len(args) - fset.NArg()
	s := args[idx]
	args = args[idx+1:]
	for _, sub := range c.subs {
		if sub.Application.Name != s {
			continue
		}

		fname := fmt.Sprintf("command `%s`", c.Application.Name)
		fs := flag.NewFlagSet(fname, sub.Application.Err)
		fs.SetOutput(out)
		handler := sub.Application.Init(fs)
		// Command specific arguments.
		if err := fs.Parse(args); err != nil {
			return err
		}
		// Command handler.
		if err := handler(args[len(args)-fs.NArg():]...); err != nil {
			return err
		}

		return run(sub, args, fs, false)
	}

	if doerror {
		return fmt.Errorf("%s is not a valid command", s)
	}

	return nil
}
