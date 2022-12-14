package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/ddkwork/golibrary/src/cpp2go/peg_cpp_parser-master/cpppeg"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// global variables
var ()

//
// utility functions
//

// Utf8Bom is BOM mark
var Utf8Bom = []byte{239, 187, 191}

func hasBOM(in []byte) bool {
	return bytes.HasPrefix(in, Utf8Bom)
}

func stripBOM(in []byte) []byte {
	return bytes.TrimPrefix(in, Utf8Bom)
}

// main routine
func main() {

	structlist := flag.Bool("struct", false, "listing for struct")
	enumlist := flag.Bool("enum", false, "listing for enum")
	debugmode := flag.Bool("d", false, "debug mode")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Printf("no input file.")
		os.Exit(1)
	}

	filename := args[0]
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if hasBOM(buf) == true {
		buf = stripBOM(buf)
	}
	parsebuffer := strings.Replace(string(buf), "\r\n", "\n", -1)
	parser := &cpppeg.Parser{Buffer: parsebuffer}
	parser.Setup(*debugmode)
	parser.Init()
	err = parser.Parse()
	if err != nil {
		fmt.Printf("%s:%s\n", filename, err)
		os.Exit(1)
	}

	parser.Execute()
	parser.Finish()

	// print function
	putfunc := func(fname string) bool {
		tpl := template.Must(template.ParseFiles(fname))
		if err := tpl.Execute(os.Stdout, parser.GetNamespace()); err != nil {
			fmt.Println(err)
			return false
		}
		return true
	}

	result := true
	if *structlist == true {
		result = putfunc("struct.tpl")
	}
	if *enumlist == true {
		result = putfunc("enum.tpl")
	}

	if result {
		fmt.Println("done.")
	} else {
		fmt.Println("error.")
		os.Exit(1)
	}

}

//
