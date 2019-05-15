package cmdflag_test

import (
	"flag"
	"fmt"
	"github.com/pierrec/cmdflag"
)

func ExampleParse() {
	// Declare the `split` cmdflag.
	c := cmdflag.New(flag.CommandLine)
	_, _ = c.Add(
		cmdflag.Application{
			Name:  "split",
			Descr: "splits files into fixed size chunks",
			Args:  "[sep ...]",
			Help: `split can split multiple files into chunks
e.g. split -size 1M file1 file2
will generate files of 1M or maybe less for the last one, as follow:
file1_0
file1_1
...
file1_x
file2_00
file2_01
...
file2_yy`,
			Err: flag.ExitOnError,
			Init: func(fs *flag.FlagSet) cmdflag.Handler {
				// Declare the cmdflag specific flags.
				var s string
				fs.StringVar(&s, "s", "", "string to be split")

				// Return the handler to be executed when the cmdflag is found.
				return func(sep ...string) (int, error) {
					i := len(s) / 2
					fmt.Printf("%s %v %s", s[:i], sep, s[i:])
					return 1, nil
				}
			},
		})

	// ./program split -s hello & @
	if err := c.Parse("split", "-s", "hello", "&", "@"); err != nil {
		panic(err)
	}

	// Output:
	// he [& @] llo
}
