package cpppeg

import (
	"fmt"
	"strconv"
)

type state int

const (
	inGlobal state = iota
	inEnum
	inStruct
)

// CppVariable is struct member information
type CppVariable struct {
	Type    string
	Name    string
	Size    string
	Comment string
}

// CppStruct is single struct information
type CppStruct struct {
	Name      string
	Comment   string
	Variables []CppVariable
}

// CppEnumValue is enum value set
type CppEnumValue struct {
	Name    string
	Value   int
	Comment string
}

// CppEnum is enum
type CppEnum struct {
	Name      string
	IsClass   bool
	ValueSize string
	Comment   string
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
	variables   []CppVariable
	enumerates  []CppEnum
	currentStr  CppStruct
	currentEnum CppEnum
	enumNumber  int
	enumSize    string
	hasNS       bool
	arraySize   string
	comment     string
	debugMode   bool
	state       state
}

//
// functions for Parser
//
func (b *Body) popLiterals(p int) {
	b.Literals = b.Literals[:len(b.Literals)-p]
}

func (b *Body) setComment(c string) {
	switch b.state {
	case inGlobal:
		b.comment = c
	case inEnum:
		top := len(b.currentEnum.EnumValue)
		if top > 0 {
			b.currentEnum.EnumValue[top-1].Comment = c
		} else {
			b.comment = c
		}
	case inStruct:
		top := len(b.variables)
		if top > 0 {
			b.variables[top-1].Comment = c
		} else {
			b.comment = c
		}
	}
}

func (b *Body) makeEnum(hasliteral bool, isclass bool) {
	name := "-"
	if hasliteral {
		stackTop := len(b.Literals)
		name = b.Literals[stackTop-1]
		b.popLiterals(1)
	}
	//fmt.Println("makeEnum:" + name)
	b.enumNumber = 0
	b.currentEnum = CppEnum{Name: name, IsClass: isclass, ValueSize: b.enumSize, EnumValue: []CppEnumValue{}, Comment: b.comment}
	b.comment = ""
	b.state = inEnum
}

func (b *Body) closeEnum() {
	b.enumerates = append(b.enumerates, b.currentEnum)
	b.comment = ""
	b.state = inGlobal
}

func (b *Body) setEnumValue() {
	stackTop := len(b.Literals)
	name := b.Literals[stackTop-1]
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
	if len(b.enumerates) > 0 {
		ns.Enumerates = append(ns.Enumerates, b.enumerates...)
	}
	return ns
}

func (b *Body) setNamespace() {
	StackTop := len(b.Literals)
	Name := b.Literals[StackTop-1]
	b.Namespace[Name] = b.addNamespace(Name)
	b.StructList = []CppStruct{}
	b.enumerates = []CppEnum{}
	b.enumSize = "int"
	b.popLiterals(1)
}

func (b *Body) makeStruct() {
	stackTop := len(b.Literals)
	name := b.Literals[stackTop-1]
	b.currentStr = CppStruct{Name: name, Comment: b.comment}
	b.comment = ""
	b.popLiterals(1)
	b.state = inStruct
}

func (b *Body) setStruct() {
	b.currentStr.Variables = b.variables
	b.StructList = append(b.StructList, b.currentStr)
	b.variables = []CppVariable{}
	b.comment = ""
	b.state = inGlobal
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
	b.variables = append(b.variables, CppVariable{Type: VarType, Name: VarName, Size: b.arraySize})
	b.popLiterals(popnum)
	b.hasNS = false
	b.arraySize = ""
}

func (b *Body) useNamespace() {
	if b.debugMode {
		StackTop := len(b.Literals)
		name := b.Literals[StackTop-1]
		fmt.Println("namespace:" + name)
	}
	b.hasNS = true
}

func (b *Body) useArray() {
	stacktop := len(b.Literals)
	b.arraySize = b.Literals[stacktop-1]
	b.popLiterals(1)
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
	b.variables = []CppVariable{}
	b.enumerates = []CppEnum{}
	b.enumSize = "int"
	b.debugMode = debug
	b.hasNS = false
	b.comment = ""
	b.arraySize = ""
	b.state = inGlobal
}

// Finish close process parser
func (b *Body) Finish() {
	if len(b.StructList) > 0 || len(b.enumerates) > 0 {
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

// GetLineNumber is get line number on error
func (b *Body) GetLineNumber() int {
	return b.line
}
