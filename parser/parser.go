package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"text/template"
)

type Lockfile struct {
	LockfileVersion int `json:"lockfileVersion"`
	Packages        map[string]struct {
		Version   string `json:"version"`
		Resolved  string `json:"resolved"`
		Integrity string `json:"integrity"`
	} `json:"packages"`
}

type Dependency struct {
	Name      string
	Version   string
	URL       string
	Integrity string
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: parser <path-to-package-lock.json>")
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var lock Lockfile
	if err := json.Unmarshal(data, &lock); err != nil {
		log.Fatal(err)
	}

	if lock.LockfileVersion < 3 {
		log.Fatal("Unsupported lockfile version")
	}

	var deps []Dependency
	for name, dep := range lock.Packages {
		if name == "" {
			continue
		}
		name := strings.TrimPrefix(name, "node_modules/")
		deps = append(deps, Dependency{name, dep.Version, dep.Resolved, dep.Integrity})
	}

	tmpl, err := template.ParseFiles("parser/derivation.gotmpl")
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.Execute(os.Stdout, deps); err != nil {
		log.Fatal(err)
	}
}
