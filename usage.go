package cmdflag

import (
	"flag"
	"fmt"
	"io"
)

// usage returns the default function used to display the help message.
func usage(c *Command) func() {
	return func() {
		out := fsetOutput(c.fset)

		_, _ = fmt.Fprintf(out, "Usage of %s:\n", program())
		c.fset.PrintDefaults()

		_, _ = fmt.Fprintf(out, "\nSubcommands:")
		for _, c := range c.Commands() {
			usageCommand(out, c.Application)
		}
	}
}

// usageCommand returns the default function used to display the help message for a given command.
func usageCommand(out io.Writer, app Application) func() {
	return func() {
		_, _ = fmt.Fprintf(out, "Usage of command `%s`:\n", app.Name)
		_, _ = fmt.Fprintf(out, "\n%s\n%s %s\n", app.Descr, app.Name, app.Args)
		fs := flag.NewFlagSet(app.Name, app.Err)
		_ = app.Init(fs)
		fs.PrintDefaults()
	}
}
