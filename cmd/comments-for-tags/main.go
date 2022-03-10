package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/stellentus/l5x"
)

var (
	in  = flag.String("in", "test.L5X", "Path to L5X file to parse")
	out = flag.String("out", "", "Path to file to write generated tag/comment CSV, or \"\" for stdout")
)

func main() {
	flag.Parse()
	var err error

	content, err := l5x.NewFromFile(*in)
	panicIfError(err, "Could not parse L5X '"+*in+"'")

	tl, err := content.Controller.TypeList()
	panicIfError(err, "Coundn't register ControlLogixTypes")

	var strB strings.Builder
	err = content.Controller.WriteTagDescriptions(tl, &strB)
	panicIfError(err, "Failed to write tags")

	fout := io.WriteCloser(os.Stdout)
	if *out != "" {
		fout, err = os.Create(*out)
		panicIfError(err, "Could not open output file '"+*out+"'")
	}
	defer fout.Close()

	_, err = fout.Write([]byte(strB.String()))
	panicIfError(err, "Failed to write output file")
}

func panicIfError(err error, reason string) {
	if err != nil {
		panic("ERROR " + err.Error() + ": " + reason)
	}
}
