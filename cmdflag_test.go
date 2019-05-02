package cmdflag_test

import (
	"flag"
	"os"
	"testing"

	"github.com/pierrec/cmdflag"
)

func TestGlobalFlagOnly(t *testing.T) {
	cmd, args := flag.CommandLine, os.Args
	defer func() { flag.CommandLine, os.Args = cmd, args }()

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
	cmd, args := flag.CommandLine, os.Args
	defer func() { flag.CommandLine, os.Args = cmd, args }()

	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)
	os.Args = []string{"program", "invalidsub"}

	if err := cmdflag.Parse(); err == nil {
		t.Fatal("expected invalid command error")
	}
}

func TestOnecmdflag(t *testing.T) {
	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
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

	args := os.Args
	defer func() { os.Args = args }()
	os.Args = []string{"program", "sub1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestOnecmdflagOneFlag(t *testing.T) {
	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
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

	args := os.Args
	defer func() { os.Args = args }()
	os.Args = []string{"program", "sub1flag", "-v1=cli1"}

	if err := cmdflag.Parse(); err != nil {
		t.Fatal(err)
	}

	if got, want := h1, 1; got != want {
		t.Fatalf("got %d; want %d", got, want)
	}
}

func TestGlobalFlagOnecmdflag(t *testing.T) {
	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
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

	cmd, args := flag.CommandLine, os.Args
	defer func() { flag.CommandLine, os.Args = cmd, args }()

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
	h1 := 0
	handle1 := func(fset *flag.FlagSet) cmdflag.Handler {
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

	cmd, args := flag.CommandLine, os.Args
	defer func() { flag.CommandLine, os.Args = cmd, args }()

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
