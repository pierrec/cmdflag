package cmdflag

import (
	"flag"
	"fmt"
)

// usage returns the default function used to display the help message.
func usage(c *Command) func() {
	return func() {
		out := fsetOutput(c.fset)

		name := c.Application.Name
		if c.Application.Init != nil {
			// Not the program.
			name = "command `" + name + "`"
		}
		_, _ = fmt.Fprintf(out, "Usage of %s:\n", name)
		c.fset.PrintDefaults()

		_, _ = fmt.Fprintf(out, "\nSubcommands:\n")
		for _, c := range c.Commands() {
			app := c.Application
			_, _ = fmt.Fprintf(out, "Usage of command `%s`:\n", app.Name)
			_, _ = fmt.Fprintf(out, "%s\n%s %s\n", app.Descr, app.Name, app.Args)
			fs := flag.NewFlagSet(app.Name, app.Err)
			_ = app.Init(fs)
			fs.PrintDefaults()
		}
	}
}
