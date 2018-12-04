package metrics

import (
	"bufio"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strings"

	"github.com/kisielk/gotool"
)

//Result - container for analysis results
type Result struct {
	LOC   int // line of code
	CLOC  int // comment line of code
	NCLOC int //

	Package int

	Struct    int
	Interface int

	Method           int
	ExportedMethod   int
	MethodLOC        int
	Function         int
	ExportedFunction int
	FunctionLOC      int

	Import int

	IfStatement     int
	SwitchStatement int
	GoStatement     int

	Test      int
	Benchmark int
	Assertion int
}

//Parser - Code parser struct
type Parser struct{}

//ParseDir - Parse all files within directory
func (p *Parser) ParsePackages(packageSpec []string) *Result {

	res := &Result{}
	for _, pac := range gotool.ImportPaths(packageSpec) {
		if pac != "." {
			//create the file set
			fset := token.NewFileSet()

			if strings.Contains(pac, "/vendor/") {
				continue
			}

			res.Package++

			fullpac := os.Getenv("GOPATH") + "/src/" + pac
			d, err := parser.ParseDir(fset, fullpac, nil, parser.ParseComments)
			if err != nil {
				log.Fatal(err)
			}

			//count up lines
			fset.Iterate(func(file *token.File) bool {
				loc, cloc, assertions := p.CountLOC(file.Name())
				res.LOC += loc
				res.CLOC += cloc
				res.NCLOC += (loc - cloc)
				res.Assertion += assertions
				return true
			})

			//setup visitors
			var visitors []AstVisitor
			visitors = append(
				visitors,
				&TypeVisitor{res: res},
				&FuncVisitor{res: res, fset: fset},
				&ImportVisitor{res: res},
				&FlowControlVisitor{res: res})

			//count entities
			for _, pkg := range d {
				ast.Inspect(pkg, func(n ast.Node) bool {
					for _, vis := range visitors {
						vis.Visit(n)
					}
					return true
				})
			}
		}
	}
	return res
}

//CountLOC - count lines of code, pull LOC, Comments, assertions
func (p *Parser) CountLOC(filePath string) (int, int, int) {

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)

	var loc int
	var cloc int
	var assertions int
	var inBlockComment bool

	assertionPrefixes := []string{
		"So(",
		"convey.So(",
		"assert.",
	}

	for {
		line, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			return loc, cloc, assertions
		}
		if isPrefix == true {
			continue //incomplete line
		}
		if len(line) == 0 {
			continue //empty line
		}
		if strings.Index(strings.TrimSpace(string(line)), "//") == 0 {
			cloc++ //slash comment at start of line
			continue
		}
		for _, prefix := range assertionPrefixes {
			if strings.HasPrefix(strings.TrimSpace(string(line)), prefix) {
				assertions++
			}
		}

		blockCommentStartPos := strings.Index(strings.TrimSpace(string(line)), "/*")
		blockCommentEndPos := strings.LastIndex(strings.TrimSpace(string(line)), "*/")

		if blockCommentStartPos > -1 {
			//block was started and not terminated
			if blockCommentEndPos == -1 || blockCommentStartPos > blockCommentEndPos {
				inBlockComment = true
			}
		}
		if blockCommentEndPos > -1 {
			//end of block is found and no new block was started
			if blockCommentStartPos == -1 || blockCommentEndPos > blockCommentStartPos {
				inBlockComment = false
				cloc++ //end of block counts as a comment line but we're already out of the block
			}
		}

		loc++
		if inBlockComment {
			cloc++
		}
	}
}

// TPC = (BT + BF + SC + MC)/(2*B + S + M) * 100%
//
// where
//
// BT - branches that evaluated to "true" at least once
// BF - branches that evaluated to "false" at least once
// SC - statements covered
// MC - methods entered
//
// B - total number of branches
// S - total number of statements
// M - total number of methods
