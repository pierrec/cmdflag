package cmdflag_test

import (
	"flag"
	"os"
	"testing"

	"github.com/pierrec/cmdflag"
)

func prepareArgs() (done func()) {
	cmd, args := flag.CommandLine, os.Args
	return func() { flag.CommandLine, os.Args = cmd, args }
}

func TestCommand_Add(t *testing.T) {
	defer prepareArgs()()

	ini := func(*flag.FlagSet) cmdflag.Initializer {
		return func(s ...string) error {
			return nil
		}
	}

	apps := func(app ...cmdflag.Application) []cmdflag.Application { return app }
	for _, tcase := range []struct {
		label string
		app   []cmdflag.Application
		err   error
	}{
		{label: "missing command", err: cmdflag.ErrMissingCommandName,
			app: apps(cmdflag.Application{})},
		{label: "missing initializer", err: cmdflag.ErrMissingInitializer,
			app: apps(cmdflag.Application{Name: "test"})},
		{label: "duplicate command", err: cmdflag.ErrDuplicateCommand,
			app: apps(cmdflag.Application{Name: "test", Init: ini}, cmdflag.Application{Name: "test", Init: ini})},
		{label: "cmd1 cmd2",
			app: apps(cmdflag.Application{Name: "cmd1", Init: ini}, cmdflag.Application{Name: "cmd2", Init: ini})},
	} {
		t.Run(tcase.label, func(t *testing.T) {
			var c cmdflag.Command
			var cmds []*cmdflag.Command
			var err error
			for _, app := range tcase.app {
				var cc *cmdflag.Command
				cc, err = c.Add(app)
				if err != nil {
					break
				}
				cmds = append(cmds, cc)
			}
			if err != nil {
				if tcase.err == nil {
					t.Fatal(err)
				}
				if got, want := err, tcase.err; got != want {
					t.Fatalf("got %#v; want %#v", got, want)
				}
				return
			}
			if tcase.err != nil {
				t.Fatal("expected error not found")
			}
			if got, want := len(c.Commands()), len(cmds); got != want {
				t.Fatalf("got %#v; want %#v", got, want)
			}
		})
	}
}

func TestGlobalFlagOnly(t *testing.T) {
	defer prepareArgs()()

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")
	os.Args = []string{"program", "-v1=gcli1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := gv1, "gcli1"; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestInvalidcmdflag(t *testing.T) {
	defer prepareArgs()()

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	os.Args = []string{"program", "invalidsub"}

	if err := cmdflag.Parse(); err == nil {
		t.Fatal("expected invalid command error")
	}
}

func TestOnecmdflag(t *testing.T) {
	defer prepareArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Initializer {
		return func(args ...string) error {
			h1++
			return nil
		}
	}
	app := cmdflag.Application{Name: "sub1", Err: flag.ExitOnError, Init: handle1}
	_, err := cmdflag.CommandLine.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"program", "sub1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestOnecmdflagOneFlag(t *testing.T) {
	defer prepareArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Initializer {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) error {
			if got, want := v1, "cli1"; got != want {
				t.Fatalf("got %s; want %s", got, want)
			}
			return nil
		}
	}
	app := cmdflag.Application{Name: "sub1flag", Err: flag.ExitOnError, Init: handle1}
	_, err := cmdflag.CommandLine.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	os.Args = []string{"program", "sub1flag", "-v1=cli1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestGlobalFlagOnecmdflag(t *testing.T) {
	defer prepareArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Initializer {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) error {
			return nil
		}
	}
	app := cmdflag.Application{Name: "subglobal", Err: flag.ExitOnError, Init: handle1}
	_, err := cmdflag.CommandLine.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")

	os.Args = []string{"program", "-v1=gcli1", "subglobal"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

	if got, want := gv1, "gcli1"; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestGlobalFlagOnecmdflagOneFlag(t *testing.T) {
	defer prepareArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Initializer {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) error {
			if got, want := v1, "cli1"; got != want {
				t.Fatalf("got %s; want %s", got, want)
			}
			return nil
		}
	}
	app := cmdflag.Application{Name: "subglobal1flag", Err: flag.ExitOnError, Init: handle1}
	_, err := cmdflag.CommandLine.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")

	os.Args = []string{"program", "subglobal1flag", "-v1=cli1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}
