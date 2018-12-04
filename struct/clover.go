package clover

import "encoding/xml"

type Coverage struct {
	XMLName   xml.Name `xml:"coverage"`
	Generated int64    `xml:"generated,attr"`
	Clover    string   `xml:"clover,attr"`
	Project   *Project `xml:"project"`
}
type Project struct {
	XMLName   xml.Name `xml:"project"`
	Timestamp int64    `xml:"timestamp,attr"`
	Name      string   `xml:"name,attr"`
	Metrics   *Metrics `xml:"metrics"`
}

type Metrics struct {
	//XMLName 		xml.Name `xml:"metrics"`
	Loc              *int `xml:"loc,omitempty,attr"`
	Cloc             *int `xml:"cloc,omitempty,attr"`
	NCLoc            *int `xml:"ncloc,omitempty,attr"`
	Files            *int `xml:"files,omitempty,attr"`
	Packages         *int `xml:"packages,omitempty,attr"`
	Struct           *int `xml:"struct,omitempty,attr"`
	Interface        *int `xml:"interface,omitempty,attr"`
	Methods          *int `xml:"methods,omitempty,attr"`
	ExportedMethod   *int `xml:"exportedmethod,omitempty,attr"`
	MethodLOC        *int `xml:"methodloc,omitempty,attr"`
	Function         *int `xml:"function,omitempty,attr"`
	ExportedFunction *int `xml:"exportedfunction,omitempty,attr"`
	FunctionLOC      *int `xml:"functionloc,omitempty,attr"`
	Import           *int `xml:"import,omitempty,attr"`
	IfStatement      *int `xml:"ifstatement,omitempty,attr"`
	SwitchStatement  *int `xml:"switchStatement,omitempty,attr"`
	GoStatement      *int `xml:"gostatement,omitempty,attr"`
	Test             *int `xml:"test,omitempty,attr"`
	Benchmark        *int `xml:"benchmark,omitempty,attr"`
	Assertion        *int `xml:"assertion,omitempty,attr"`
	Statements       *int `xml:"statements,omitempty,attr"`
	CoverStatements  *int `xml:"coveredstatements,omitempty,attr"`
	CoverMethods     *int `xml:"coveredmethods,omitempty,attr"`
	Elements         *int `xml:"elements,omitempty,attr"`
	Coveredelements  *int `xml:"coveredelements,omitempty,attr"`

	Package []*Package `xml:"package"`
}

type Package struct {
	XMLName xml.Name `xml:"package"`
	Name    string   `xml:"name,attr"`
	Metrics *Metric  `xml:"metrics"`
	Files   []*File  `xml:"file"`
}

type Metric struct {
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

type Line struct {
	XMLName    xml.Name `xml:"line"`
	Num        *int     `xml:"num,attr"`
	Count      *int     `xml:"count,attr"`
	Type       *string  `xml:"type,attr"`
	Truecount  *string  `xml:"truecount,attr"`
	Falsecount *string  `xml:"falsecount,attr"`
}

type File struct {
	XMLName xml.Name `xml:"file"`
	Name    string   `xml:"name,attr"`
	Path    *string  `xml:"path,attr"`
	Metrics *Metric  `xml:"metrics"`
	Lines   []*Line  `xml:"line,chardata`
}
