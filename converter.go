package main

import (
	"path"
	"time"

	"github.com/codeofthrone/goclover/cover"
	"github.com/codeofthrone/goclover/metrics"
)

var packageMemo = make(map[string]*Package)
var fileMemo = make(map[string]*File)
var lineMemo = make(map[string]map[int]bool)
var stmt = "stmt"

// ConvertProfilesToClover
func ConvertProfilesToClover(profs []*cover.Profile) (*Coverage, error) {
	for i := range profs {
		prof := profs[i]
		fname := prof.FileName
		d, _ := path.Split(fname)
		packagename := path.Dir(d)
		p, _ := findOrCreatePackage(packagename)
		f, created := findOrCreateFile(fname)
		if created {
			p.Files = append(p.Files, f)
		}
		lines, ok := lineMemo[fname]
		if !ok {
			lines = make(map[int]bool)
		}
		for j := range prof.Blocks {
			block := prof.Blocks[j]
			for k := block.StartLine; k < block.EndLine; k++ {
				if _, already := lines[k]; already {
					continue
				}
				num := k
				count := block.Count
				t := stmt
				l := &Line{
					Num:   &num,
					Count: &count,
					Type:  &t,
				}
				lines[k] = true
				f.Lines = append(f.Lines, l)
			}
		}
		lineMemo[fname] = lines
	}
	pr := &Project{}
	pr.Timestamp = time.Now().Unix()
	packages := make([]*Package, len(packageMemo))
	i := 0
	for k := range packageMemo {
		mp := metrics.Parser{}
		var pacName []string
		pacName = append(pacName, packageMemo[k].Name)
		res := mp.ParsePackages(pacName)
		p := &Package{
			XMLName: packageMemo[k].XMLName,
			Name:    packageMemo[k].Name,
			Metrics: &Metrics{
				Loc:              &res.LOC,
				NCLoc:            &res.NCLOC,
				Cloc:             &res.CLOC,
				Packages:         &res.Package,
				Struct:           &res.Struct,
				Method:           &res.Method,
				ExportedMethod:   &res.ExportedMethod,
				MethodLOC:        &res.MethodLOC,
				Function:         &res.Package,
				ExportedFunction: &res.ExportedFunction,
				FunctionLOC:      &res.FunctionLOC,
				Import:           &res.Import,
				IfStatement:      &res.IfStatement,
				SwitchStatement:  &res.SwitchStatement,
				Test:             &res.Test,
				Benchmark:        &res.Benchmark,
				Assertion:        &res.Assertion,
			},
			Files: packageMemo[k].Files,
		}
		packages[i] = p
		i++
	}
	pr.Packages = packages
	c := &Coverage{}
	c.GeneratedAt = time.Now().Unix()
	c.AppProject = pr
	return c, nil
}

func findOrCreateFile(name string) (*File, bool) {
	f, ok := fileMemo[name]
	if ok {
		return f, false
	}
	f = &File{}
	f.Name = name
	fileMemo[name] = f
	f.Lines = make([]*Line, 0)
	f.Classes = make([]*Class, 0)
	return f, true
}

func findOrCreatePackage(name string) (*Package, bool) {
	p, ok := packageMemo[name]
	if ok {
		return p, false
	}
	p = &Package{}
	p.Name = name
	packageMemo[name] = p
	p.Files = make([]*File, 0)
	return p, true
}
