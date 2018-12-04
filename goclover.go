package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/codeofthrone/goclover/cover"
	"github.com/codeofthrone/goclover/metrics"
	"github.com/codeofthrone/goclover/struct"
)

var packageMemo = make(map[string]*clover.Package)
var fileMemo = make(map[string]*clover.File)
var lineMemo = make(map[string]map[int]bool)
var stmt = "stmt"
var coverageP float64

func findOrCreateFile(name string) (*clover.File, bool) {
	f, ok := fileMemo[name]
	if ok {
		return f, false
	}
	f = &clover.File{}
	f.Name = name
	fileMemo[name] = f
	f.Lines = make([]*clover.Line, 0)
	return f, true
}

func findOrCreatePackage(name string) (*clover.Package, bool) {
	p, ok := packageMemo[name]
	if ok {
		return p, false
	}
	p = &clover.Package{}
	p.Name = name
	packageMemo[name] = p
	p.Files = make([]*clover.File, 0)
	return p, true
}

func findOrCreateMetrics() (*clover.Metrics, bool) {
	p := &clover.Metrics{}
	p.Package = make([]*clover.Package, 0)
	return p, true
}

func StructToJsonString(input interface{}) string {
	byteinput, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
	}
	return string(byteinput)
}

// ConvertProfilesToClover
func ConvertProfilesToClover(profs []*cover.Profile) (*clover.Coverage, error) {
	filecount := 0
	for i := range profs {
		prof := profs[i]
		fname := prof.FileName
		filecount++
		d, _ := path.Split(fname)
		packagename := path.Dir(d)
		m, _ := findOrCreateMetrics()
		p, _ := findOrCreatePackage(packagename)
		f, created := findOrCreateFile(fname)
		if created {
			p.Files = append(p.Files, f)
		}
		m.Package = append(m.Package, p)
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
				l := &clover.Line{
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
	pr := &clover.Project{}
	pr.Timestamp = time.Now().Unix()

	totalpac := []string{}
	totalpackage := []*clover.Package{}
	for k := range packageMemo {
		totalpac = append(totalpac, k)

		mp := metrics.Parser{}
		res := mp.ParsePackages([]string{k})
		packageMemo[k].Metrics = &clover.Metric{
			Loc:      &res.LOC,
			NCLoc:    &res.NCLOC,
			Cloc:     &res.CLOC,
			Packages: &res.Package,
		}
		totalpackage = append(totalpackage, packageMemo[k])
	}

	var loc, cloc, ncloc, packages int
	mp := metrics.Parser{}
	res := mp.ParsePackages(totalpac)
	loc += res.LOC
	cloc += res.CLOC
	ncloc += res.NCLOC
	packages += res.Package

	tmpvalue := 1000
	tmpres := int(coverageP * 10)
	pr.Metrics = &clover.Metrics{
		Loc:             &loc,
		Cloc:            &cloc,
		NCLoc:           &ncloc,
		Packages:        &packages,
		Package:         totalpackage,
		Files:           &filecount,
		Methods:         &tmpvalue,
		CoverMethods:    &tmpres,
		Statements:      &tmpvalue,
		CoverStatements: &tmpres,
		Elements:        &tmpvalue,
		Coveredelements: &tmpres,
	}
	c := &clover.Coverage{}
	c.Clover = "3.2.0"
	c.Project = pr
	c.Project.Name = "All files"
	c.Generated = time.Now().Unix()
	return c, nil
}

func consume(wg *sync.WaitGroup, r io.Reader) {
	defer wg.Done()
	reader := bufio.NewReader(r)
	for {
		l, _, err := reader.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		fmt.Println(string(l))
	}
}

func gotest(args []string) int {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	r, w := io.Pipe()
	defer w.Close()

	args = append([]string{"test"}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stderr = w
	cmd.Stdout = w
	cmd.Env = os.Environ()

	go consume(&wg, r)

	if err := cmd.Run(); err != nil {
		if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus()
		}
		return 1
	}
	return 0
}

func gotool(args []string) {
	cmd := exec.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(opBytes), "\n")
	for _, value := range lines {
		if strings.Contains(value, "total:") {
			rxp := regexp.MustCompile(`\(statements\).*\t(.*)%`)
			res := rxp.FindStringSubmatch(value)
			if len(res) == 2 {
				tmpfloat, _ := strconv.ParseFloat(res[1], 1)
				coverageP = tmpfloat
			}
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	var fp, rp string
	flag.StringVar(&fp, "f", "c.out", "coverprofile file name")
	flag.StringVar(&rp, "o", "clover.xml", "clover XML save to file name")
	flag.Parse()

	testcmds := []string{"-v", "-coverprofile=" + fp, "-covermode=atomic", "./..."}
	totalcmds := []string{"tool", "cover", "-func=" + fp}
	gotest(testcmds)
	gotool(totalcmds)

	files, err := filepath.Glob(fp)
	if err != nil {
		panic(err)
	}
	p := make([]*cover.Profile, 0)
	for _, f := range files {
		tmpP, err := cover.ParseProfiles(f)
		if err != nil {
			panic(err)
		}
		p = append(p, tmpP...)
	}
	c, err := ConvertProfilesToClover(p)
	if err != nil {
		log.Println(err)
	}

	out, err := os.OpenFile(rp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Println(err)
	}

	defer out.Close()
	fmt.Fprintf(out, xml.Header)
	encoder := xml.NewEncoder(out)
	encoder.Indent("", "\t")
	if err = encoder.Encode(c); err != nil {
		panic(err)
	}
	log.Println("Coverage :", coverageP, "%")
}
