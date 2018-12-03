package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/codeofthrone/goclover/cover"
)

func main() {
	log.SetFlags(log.Lshortfile)

	var fp, rp string
	flag.StringVar(&fp, "file", "*.coverprofile", "file path string or file path regex string")
	flag.StringVar(&fp, "f", "c.out", "file path string or file path regex string")
	flag.StringVar(&rp, "o", "clover.xml", "clover XML save to file name")
	flag.Parse()
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
	c, _ := ConvertProfilesToClover(p)

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

}
