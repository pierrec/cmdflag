// +build go1.12

package cmdflag

import (
	"fmt"
	"runtime/debug"
	"strings"
)

func buildinfo() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Version
	}
	return "no version available (not built with module support)"
}

func fullbuildinfo() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "no version available (not built with module support)"
	}
	var buf strings.Builder
	printModule(&buf, &bi.Main)
	for _, m := range bi.Deps {
		buf.WriteByte('\t')
		printModule(&buf, m)
	}
	return buf.String()
}

func printModule(buf *strings.Builder, m *debug.Module) {
	_, _ = fmt.Fprintf(buf, "%s %s", m.Path, m.Version)
	if m.Replace != nil {
		buf.WriteString(" => ")
		printModule(buf, m.Replace)
	} else {
		buf.WriteByte('\n')
	}
}
