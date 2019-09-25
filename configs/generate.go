//+build ignore
//go:generate -command asset go run generate.go
//go:generate asset policies.yaml

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type asset struct {
	Name        string
	Description string
	Content     string
}

var (
	desc = flag.String("desc", "auto-generated from asset file", "description")
)

func main() {
	flag.Parse()
	log.Println("Starting Asset Generation")

	for _, filename := range flag.Args() {
		log.Println("Generate from", filename)

		dir, base := filepath.Split(filename)

		if dir == "" {
			dir = "."
		}

		asset := asset{}
		asset.Name = strings.Title(strings.TrimSuffix(base, filepath.Ext(base)))
		asset.Description = *desc

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		asset.Content = string(content)

		t, err := template.ParseFiles("config.tmpl")
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		f, err := os.Create(base + ".go")
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		err = t.Execute(f, asset)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		f.Close()
	}
}
