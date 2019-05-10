// +build go1.12

package cmdflag

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const nobuildinfo = "no version available (not built with module support)"

func buildinfo() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Version
	}
	return nobuildinfo
}

func fullbuildinfo() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return nobuildinfo
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
