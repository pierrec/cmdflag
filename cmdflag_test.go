package cmdflag_test

import (
	"bytes"
	"flag"
	"os"
	"testing"

	"github.com/pierrec/cmdflag"
)

func restoreArgs() (done func()) {
	cmd, args := flag.CommandLine, os.Args
	flag.CommandLine = flag.NewFlagSet("test", 0)
	return func() { flag.CommandLine, os.Args = cmd, args }
}

func TestCommand_Add(t *testing.T) {
	defer restoreArgs()()

	ini := func(*flag.FlagSet) cmdflag.Handler {
		return func(s ...string) (int, error) {
			return 0, nil
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
			c := cmdflag.New(nil)
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

func TestCommand_MustAdd(t *testing.T) {
	defer restoreArgs()()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on invalid application")
		}
	}()

	c := cmdflag.New(nil)
	c.MustAdd(cmdflag.Application{})
}

func TestGlobalFlagOnly(t *testing.T) {
	defer restoreArgs()()

	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")

	c := cmdflag.New(nil)
	if err := c.Parse("-v1=gcli1"); err != nil {
		t.Fatal(err)
	}

	if got, want := gv1, "gcli1"; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestInvalidCommand(t *testing.T) {
	defer restoreArgs()()

	c := cmdflag.New(nil)
	app := cmdflag.Application{
		Name: "sub1",
		Init: func(fset *flag.FlagSet) cmdflag.Handler { return nil },
	}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Parse("invalidsub"); err == nil {
		t.Fatal("expected invalid command error")
	}
}

func TestNoCommandSet(t *testing.T) {
	defer restoreArgs()()

	c := cmdflag.New(nil)
	if err := c.Parse("sub"); err != nil {
		t.Fatal(err)
	}
}

func TestOneCommand(t *testing.T) {
	h := 0
	handle := func(fset *flag.FlagSet) cmdflag.Handler {
		return func(args ...string) (int, error) {
			h++
			return 0, nil
		}
	}
	c := cmdflag.New(nil)
	app := cmdflag.Application{Name: "sub1", Err: flag.ExitOnError, Init: handle}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("sub1"); err != nil {
		t.Fatal(err)
	}

	if got, want := h, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestOneCommandOneNestedCommand(t *testing.T) {
	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
		return func(args ...string) (int, error) {
			h1++
			return 0, nil
		}
	}
	c1 := cmdflag.New(nil)
	app1 := cmdflag.Application{Name: "sub1", Err: flag.ExitOnError, Init: handle1}
	c2, err := c1.Add(app1)
	if err != nil {
		t.Fatal(err)
	}
	h2 := 0
	handle2 := func(fset *flag.FlagSet) cmdflag.Handler {
		return func(args ...string) (int, error) {
			h2 += 10
			return 0, nil
		}
	}
	app2 := cmdflag.Application{Name: "sub2", Err: flag.ExitOnError, Init: handle2}
	_, err = c2.Add(app2)
	if err != nil {
		t.Fatal(err)
	}

	if err := c1.Parse("sub1", "sub2"); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
	if got, want := h2, 10; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestOneCommandOneFlag(t *testing.T) {
	defer restoreArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) (int, error) {
			if got, want := v1, "cli1"; got != want {
				t.Fatalf("got %s; want %s", got, want)
			}
			return 0, nil
		}
	}
	c := cmdflag.New(nil)
	app := cmdflag.Application{Name: "sub1flag", Err: flag.ExitOnError, Init: handle1}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("sub1flag", "-v1=cli1"); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestGlobalFlagOneCommand(t *testing.T) {
	defer restoreArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) (int, error) {
			return 0, nil
		}
	}
	c := cmdflag.New(nil)
	app := cmdflag.Application{Name: "subglobal", Err: flag.ExitOnError, Init: handle1}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")

	if err := c.Parse("-v1=gcli1", "subglobal"); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}

	if got, want := gv1, "gcli1"; got != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestGlobalFlagOneCommandOneFlag(t *testing.T) {
	defer restoreArgs()()

	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
		h1++

		var v1 string
		fset.StringVar(&v1, "v1", "val1", "usage1")

		return func(args ...string) (int, error) {
			if got, want := v1, "cli1"; got != want {
				t.Fatalf("got %s; want %s", got, want)
			}
			return 0, nil
		}
	}
	c := cmdflag.New(nil)
	app := cmdflag.Application{Name: "subglobal1flag", Err: flag.ExitOnError, Init: handle1}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	var gv1 string
	flag.StringVar(&gv1, "v1", "val1", "usage1")

	if err := c.Parse("subglobal1flag", "-v1=cli1"); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestVersion(t *testing.T) {
	defer restoreArgs()()

	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.CommandLine.Bool(cmdflag.VersionBoolFlag, false, "print the program version")

	c := cmdflag.New(nil)
	app := cmdflag.Application{
		Name: "sub",
		Err:  flag.ExitOnError,
		Init: func(fset *flag.FlagSet) cmdflag.Handler { return nil },
	}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("-"+cmdflag.VersionBoolFlag+"=false", "dummy"); err != cmdflag.ErrNoCommand {
		t.Fatal("disabled version flag should error")
	}
	if err := c.Parse("-" + cmdflag.VersionBoolFlag); err != nil {
		t.Fatal(err)
	}

	if got := buf.Bytes(); !bytes.Contains(got, []byte("version")) {
		t.Fatal("version flag does not output the version")
	}
}

func TestFullVersion(t *testing.T) {
	defer restoreArgs()()

	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)
	flag.CommandLine.Bool(cmdflag.FullVersionBoolFlag, false, "print the program version")

	c := cmdflag.New(nil)
	app := cmdflag.Application{
		Name: "sub",
		Err:  flag.ExitOnError,
		Init: func(fset *flag.FlagSet) cmdflag.Handler { return nil },
	}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("-"+cmdflag.FullVersionBoolFlag+"=false", "dummy"); err != cmdflag.ErrNoCommand {
		t.Fatal("disabled full version flag should error")
	}

	if err := c.Parse("-" + cmdflag.FullVersionBoolFlag); err != nil {
		t.Fatal(err)
	}

	if got := buf.Bytes(); !bytes.Contains(got, []byte("full version")) {
		t.Fatal("full version flag does not output the full version")
	}
}

func TestHelpCommand(t *testing.T) {
	defer restoreArgs()()

	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)

	c := cmdflag.New(nil)
	if err := c.AddHelp(); err != nil {
		t.Fatal(err)
	}
	app := cmdflag.Application{
		Name: "sub",
		Help: "helpcommand",
		Err:  flag.ExitOnError,
		Init: func(fset *flag.FlagSet) cmdflag.Handler { return nil },
	}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("help"); err != nil {
		t.Fatal(err)
	}
	if got := buf.Bytes(); !bytes.Contains(got, []byte("Usage")) {
		t.Fatal("help command with no command does not output the usage")
	}
	buf.Truncate(0)

	if err := c.Parse("help", "dummy"); err == nil {
		t.Fatal("help command should fail on invalid command")
	}

	if err := c.Parse("help", app.Name); err != nil {
		t.Fatal(err)
	}

	if got := buf.Bytes(); !bytes.Contains(got, []byte(app.Help)) {
		t.Fatal("full version flag does not output the full version")
	}
}

func TestHelp(t *testing.T) {
	defer restoreArgs()()

	buf := new(bytes.Buffer)
	flag.CommandLine.SetOutput(buf)

	c := cmdflag.New(nil)
	app := cmdflag.Application{
		Name: "sub",
		Err:  flag.ContinueOnError,
		Init: func(fset *flag.FlagSet) cmdflag.Handler { return nil },
	}
	_, err := c.Add(app)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Parse("-h"); err != flag.ErrHelp {
		t.Fatalf("got %v; want %v", err, flag.ErrHelp)
	}
	if got := buf.Bytes(); !bytes.Contains(got, []byte("Usage")) {
		t.Fatal("invalid usage")
	}
	buf.Truncate(0)

	if err := c.Parse(app.Name, "-h"); err != flag.ErrHelp {
		t.Fatalf("got %v; want %v", err, flag.ErrHelp)
	}

	if got := buf.Bytes(); !bytes.Contains(got, []byte("Usage of command")) {
		t.Fatal("invalid command usage")
	}
}
