package main

import (
	"encoding/xml"
)

// Coverage root element "coverage"
type Coverage struct {
	XMLName     xml.Name `xml:"coverage"`
	Version     string   `xml:"clover,attr"`
	GeneratedAt int64    `xml:"generated,attr"`
	AppProject  *Project `xml:"project"`
	TestProject *Project `xml:"testproject"`
}

// Project Project element "project" or "testproject"
type Project struct {
	Name      string `xml:"name,attr"`
	Timestamp int64  `xml:"timestamp,attr"`
	Metrics   *Metrics
	Packages  []*Package
}

// Package element "package"
type Package struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Metrics *Metrics
	Files   []*File
}

// Metrics element "metrics"
// modify for golang base on github.com/warmans/golocc
type Metrics struct {
	XMLName          xml.Name `xml:"metrics"`
	Loc              *int     `xml:"loc,omitempty,attr"`
	Cloc             *int     `xml:"cloc,omitempty,attr"`
	NCLoc            *int     `xml:"ncloc,omitempty,attr"`
	Files            *int     `xml:"files,omitempty,attr"`
	Packages         *int     `xml:"packages,omitempty,attr"`
	Struct           *int     `xml:"struct,omitempty,attr"`
	Interface        *int     `xml:"interface,omitempty,attr"`
	Method           *int     `xml:"method,omitempty,attr"`
	ExportedMethod   *int     `xml:"exportedmethod,omitempty,attr"`
	MethodLOC        *int     `xml:"methodloc,omitempty,attr"`
	Function         *int     `xml:"function,omitempty,attr"`
	ExportedFunction *int     `xml:"exportedfunction,omitempty,attr"`
	FunctionLOC      *int     `xml:"functionloc,omitempty,attr"`
	Import           *int     `xml:"import,omitempty,attr"`
	IfStatement      *int     `xml:"ifstatement,omitempty,attr"`
	SwitchStatement  *int     `xml:"switchStatement,omitempty,attr"`
	GoStatement      *int     `xml:"gostatement,omitempty,attr"`
	Test             *int     `xml:"test,omitempty,attr"`
	Benchmark        *int     `xml:"benchmark,omitempty,attr"`
	Assertion        *int     `xml:"assertion,omitempty,attr"`
}

// File file element "file"
type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Path    string   `xml:"path,attr"`
	Metrics *Metrics
	Classes []*Class
	Lines   []*Line
}

// Class file element "class"
type Class struct {
	XMLName xml.Name `xml:"class"`
	Name    string   `xml:"name,attr"`
}

// Line line element "line"
type Line struct {
	XMLName      xml.Name `xml:"line"`
	Num          *int     `xml:"num,omitempty,attr"`
	Type         *string  `xml:"type,omitempty,attr"`
	Visibility   *string  `xml:"visibility,omitempty,attr"`
	Count        *int     `xml:"count,omitempty,attr"`
	Complexity   *int     `xml:"complexity,omitempty,attr"`
	MethodSig    *string  `xml:"methodsig,omitempty,attr"`
	TestSuccess  *int     `xml:"testsuccess,omitempty,attr"`
	TestDuration *int     `xml:"testduration,omitempty,attr"`
	TrueCount    *int     `xml:"truecount,omitempty,attr"`
	FalseCount   *int     `xml:"falsecount,omitempty,attr"`
}
