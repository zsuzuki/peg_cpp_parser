package cpppeg

import (
	"fmt"
	"strconv"
)

// CppVariable is struct member information
type CppVariable struct {
	Type string
	Name string
}

// CppStruct is single struct information
type CppStruct struct {
	Name      string
	Variables []CppVariable
}

// CppEnumValue is enum value set
type CppEnumValue struct {
	Name  string
	Value int
}

// CppEnum is enum
type CppEnum struct {
	Name      string
	IsClass   bool
	ValueSize string
	EnumValue []CppEnumValue
}

// CppNamespace is data set in namespace
type CppNamespace struct {
	StructList []CppStruct
	Enumerates []CppEnum
}

// Body is parser
type Body struct {
	Literals    []string
	line        int
	Namespace   map[string]CppNamespace
	StructList  []CppStruct
	Variables   []CppVariable
	Enumerates  []CppEnum
	currentEnum CppEnum
	enumNumber  int
	enumSize    string
	hasNS       bool
	debugMode   bool
}

//
// functions for Parser
//
func (b *Body) popLiterals(p int) {
	b.Literals = b.Literals[:len(b.Literals)-p]
}

func (b *Body) makeEnum(hasliteral bool, isclass bool) {
	name := "-"
	if hasliteral {
		StackTop := len(b.Literals)
		name = b.Literals[StackTop-1]
		b.popLiterals(1)
	}
	//fmt.Println("makeEnum:" + name)
	b.enumNumber = 0
	b.currentEnum = CppEnum{Name: name, IsClass: isclass, ValueSize: b.enumSize, EnumValue: []CppEnumValue{}}
}

func (b *Body) closeEnum() {
	b.Enumerates = append(b.Enumerates, b.currentEnum)
}

func (b *Body) setEnumValue() {
	StackTop := len(b.Literals)
	name := b.Literals[StackTop-1]
	//fmt.Printf("setEnumValue: %s(%d)\n", name, b.currentEnum.Current)
	b.currentEnum.EnumValue = append(b.currentEnum.EnumValue, CppEnumValue{Name: name, Value: b.enumNumber})
	b.enumNumber++
	b.popLiterals(1)
}

func (b *Body) resetEnum(nstr string) {
	si, _ := strconv.Atoi(nstr)
	b.enumNumber = si
	//fmt.Printf("value: %d\n", si)
}

func (b *Body) setEnumSize() {
	StackTop := len(b.Literals)
	b.enumSize = b.Literals[StackTop-1]
	b.popLiterals(1)
}

func (b *Body) pushStructName(n string) {
	if b.debugMode {
		fmt.Println(n)
	}
}

func (b *Body) dump(msg string) {
	if b.debugMode {
		fmt.Println("Debug:" + msg)
	}
}

func (b *Body) addNamespace(name string) CppNamespace {
	ns, ok := b.Namespace[name]
	if ok == false {
		ns = CppNamespace{StructList: []CppStruct{}}
	}
	if len(b.StructList) > 0 {
		ns.StructList = append(ns.StructList, b.StructList...)
	}
	if len(b.Enumerates) > 0 {
		ns.Enumerates = append(ns.Enumerates, b.Enumerates...)
	}
	return ns
}

func (b *Body) setNamespace() {
	StackTop := len(b.Literals)
	Name := b.Literals[StackTop-1]
	b.Namespace[Name] = b.addNamespace(Name)
	b.StructList = []CppStruct{}
	b.Enumerates = []CppEnum{}
	b.enumSize = "int"
	b.popLiterals(1)
}

func (b *Body) setStruct() {
	StackTop := len(b.Literals)
	Name := b.Literals[StackTop-1]
	b.StructList = append(b.StructList, CppStruct{Name: Name, Variables: b.Variables})
	b.Variables = []CppVariable{}
	b.popLiterals(1)
}

func (b *Body) setVar() {
	StackTop := len(b.Literals)
	VarName := b.Literals[StackTop-1]
	VarType := b.Literals[StackTop-2]
	var popnum = 2
	if b.hasNS {
		VarType = b.Literals[StackTop-3] + "::" + VarType
		popnum++
	}
	b.Variables = append(b.Variables, CppVariable{Type: VarType, Name: VarName})
	b.popLiterals(popnum)
	b.hasNS = false
}

func (b *Body) useNamespace() {
	if b.debugMode {
		StackTop := len(b.Literals)
		name := b.Literals[StackTop-1]
		fmt.Println("namespace:" + name)
	}
	b.hasNS = true
}

func (b *Body) pushLiteral(l string) {
	b.Literals = append(b.Literals, l)
}

func (b *Body) addline() {
	b.line++
}

// Setup setup parser variables
func (b *Body) Setup(debug bool) {
	b.Literals = []string{}
	b.line = 1
	b.Namespace = map[string]CppNamespace{}
	b.StructList = []CppStruct{}
	b.Variables = []CppVariable{}
	b.Enumerates = []CppEnum{}
	b.enumSize = "int"
	b.debugMode = debug
}

// Finish close process parser
func (b *Body) Finish() {
	if len(b.StructList) > 0 || len(b.Enumerates) > 0 {
		if b.debugMode {
			fmt.Println("add global namespace(\"--\")")
		}
		b.Namespace["--"] = b.addNamespace("--")
	}
}

// GetNamespace return struct/enum list in namespace
func (b *Body) GetNamespace() map[string]CppNamespace {
	return b.Namespace
}

// GetLine is get line number on error
func (b *Body) GetLineNumber() int {
	return b.line
}
