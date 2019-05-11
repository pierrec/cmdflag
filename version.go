package cmdflag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// VersionBoolFlag is the flag name to be used as a boolean flag to display the program version.
	// Declaring a boolean flag with that name will automatically implement version display.
	VersionBoolFlag = "version"
	// FullVersionBoolFlag is the flag name to be used as a boolean flag to display the full program version,
	// Declaring a boolean flag with that name will automatically implement displaying the program version,
	// including its modules and compiler versions.
	FullVersionBoolFlag = "fullversion"
)

// hasBoolFlag returns whether or not the flag with name `name` was defined and set to true.
func hasBoolFlag(fset *flag.FlagSet, name string) bool {
	f := fset.Lookup(name)
	if f == nil {
		return false
	}
	// All values implemented by the flag package implement the flag.Getter interface.
	v, ok := f.Value.(flag.Getter)
	if !ok {
		return false
	}
	b, ok := v.Get().(bool)
	// Was the flag defined as a bool and is it set?
	return ok && b
}

// program returns a deterministic name for the running program.
func program() string {
	return strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")
}

// version prints out the short version of the running program.
func version(out io.Writer) {
	_, _ = fmt.Fprintf(out, "%s version %s %s/%s\n",
		program(), buildinfo(),
		runtime.GOOS, runtime.GOARCH)
}

// fullversion prints out the full version of the running program with compiler and modules info.
func fullversion(out io.Writer) {
	_, _ = fmt.Fprintf(out, "%s full version %s %s/%s compiled by %s (%s)\n%s\n",
		program(), buildinfo(),
		runtime.GOOS, runtime.GOARCH,
		runtime.Compiler, runtime.Version(),
		fullbuildinfo())
}
