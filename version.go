package cmdflag

import (
	"flag"
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
