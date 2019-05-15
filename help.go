package cmdflag

import (
	"flag"
	"fmt"
)

// HelpCommand is the command name used to display the help of a given command.
// If no command is supplied, the usage is displayed.
//
// To display the help of a command (Application.Help), do:
//   ./myprogram help commandname
const HelpCommand = "help"

// addHelpCommand adds the `help` command to the Command c.
func addHelpCommand(c *Command) error {
	app := Application{
		Name:  HelpCommand,
		Descr: "display the help for a given command",
		Args:  "command",
		Init: func(set *flag.FlagSet) Handler {
			return func(args ...string) (int, error) {
				if len(args) == 0 {
					c.fset.Usage()
					return 0, nil
				}
				name := args[0]
				out := fsetOutput(set)
				for _, sub := range c.subs {
					if sub.Application.Name != name {
						continue
					}
					app := sub.Application
					_, _ = fmt.Fprintf(out, "%s\n%s %s\n%s\n", app.Descr, app.Name, app.Args, app.Help)
					return 1, nil
				}
				return 1, fmt.Errorf("command %s not found", name)
			}
		},
	}
	_, err := c.Add(app)
	return err
}
