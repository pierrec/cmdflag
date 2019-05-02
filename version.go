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

func hasBoolFlag(fset *flag.FlagSet, name string) bool {
	if f := fset.Lookup(name); f != nil {
		if v, ok := f.Value.(flag.Getter); ok {
			// All values implemented by the flag package implement the flag.Getter interface.
			b, ok := v.Get().(bool)
			// Was the flag was defined as a bool and is set?
			return ok && b
		}
	}
	return false
}
