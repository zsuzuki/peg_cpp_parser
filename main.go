package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"flag"

	"./cpppeg"
)

//
// global variables
//
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

//
// main routine
//
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

	for nn, ns := range parser.GetNamespace() {
		fmt.Println("Namespace: " + nn)
		if *structlist == true {
			for _, st := range ns.StructList {
				fmt.Printf("  %s in %d comment:%s\n", st.Name, len(st.Variables), st.Comment)
				for _, sv := range st.Variables {
					fmt.Printf("    %s\t%s%s comment:%s\n", sv.Type, sv.Name, func() string {
						if sv.Size != "" {
							return "[" + sv.Size + "]"
						}
						return ""
					}(), sv.Comment)
				}
			}
		}
		if *enumlist == true {
			for _, enum := range ns.Enumerates {
				fmt.Printf("Enum[%s]: %s%s comment:%s\n", enum.Name, enum.ValueSize, func() string {
					if enum.IsClass {
						return "(class)"
					}
					return ""
				}(), enum.Comment)
				for _, ev := range enum.EnumValue {
					fmt.Printf("  %s = %d comment:%s\n", ev.Name, ev.Value, ev.Comment)
				}
			}
		}
	}
	fmt.Println("done.")
}

//
