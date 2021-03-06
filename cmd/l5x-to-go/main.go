package main

import (
	"flag"
	"go/format"
	"io"
	"os"
	"strings"
	"time"

	"github.com/stellentus/l5x"
)

var (
	in  = flag.String("in", "test.L5X", "Path to L5X file to parse")
	out = flag.String("out", "", "Path to file to write generated types, or \"\" for stdout")
	pkg = flag.String("package", "", "If set, this string will be used as the package name. Otherwise, no package name will be printed.")
)

func main() {
	flag.Parse()
	var err error

	content, err := l5x.NewFromFile(*in)
	panicIfError(err, "Could not parse L5X '"+*in+"'")

	tl, err := content.Controller.TypeList()
	panicIfError(err, "Coundn't register ControlLogixTypes")

	var src strings.Builder

	// Print header line if requested
	if *pkg != "" {
		src.WriteString("// Code generated by l5x-to-go. DO NOT EDIT.\n")
		src.WriteString("// Generated from " + *in + "\n")
		src.WriteString("// Generated on " + time.Now().Format(time.RFC3339) + "\n")
		src.WriteString("package " + *pkg + "\n\n")
	}

	err = tl.WriteDefinitions(&src)
	panicIfError(err, "Failed to write definitions")

	err = content.Controller.WriteTagsStruct(tl, &src)
	panicIfError(err, "Failed to write tags")

	str, err := format.Source([]byte(src.String()))
	panicIfError(err, "Failed to run go format")

	fout := io.WriteCloser(os.Stdout)
	if *out != "" {
		fout, err = os.Create(*out)
		panicIfError(err, "Could not open output file '"+*out+"'")
	}
	defer fout.Close()

	_, err = fout.Write(str)
	panicIfError(err, "Failed to write output file")
}

func panicIfError(err error, reason string) {
	if err != nil {
		panic("ERROR " + err.Error() + ": " + reason)
	}
}
