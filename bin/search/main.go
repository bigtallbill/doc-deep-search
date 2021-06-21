package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gitlab.com/bigtallbill/doc-deep-search/pkg/docx"
	"gitlab.com/bigtallbill/doc-deep-search/pkg/xlsx"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	directory = kingpin.Arg("dir", "the directory to search in").Required().String()
	term      = kingpin.Arg("term", "the search term").Required().String()
	makeHtml  = kingpin.Flag("out", "generates a html file of results and opens it").Short('o').Bool()
	foundList []Find
)

func main() {
	kingpin.Parse()

	err := filepath.Walk(*directory, visit)

	if err != nil {
		log.Panicf("cannot walk dir '%s': %s", *directory, err)
	}

	fmt.Printf("Found %s in %d documents:\n\n", *term, len(foundList))

	for _, find := range foundList {
		fmt.Printf("%dx in \"%s\"\n", find.Occurrences, find.Path)
	}

	if *makeHtml {
		var t = template.New("home.html") // Create a template.
		t, err = t.Parse(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
	   <meta charset="UTF-8">
	   <title>Search Results for {{.Term}}</title>
		<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" integrity="sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB" crossorigin="anonymous">
	</head>
	<body>
	<div class="container">
	<h2>"{{.Term}}" Results:</h2>
		<table class="table">
	<thead>
	<tr>
	     <th scope="col">Count</th>
	     <th scope="col">File</th>
	   </tr>
	</thead>
		<tbody>
	{{range .Results}}
	<tr>
	<td>{{.Occurrences}}</td>
	<td><a href="file:///{{.AbsPath}}">{{.Path}}</a></td>
	</tr>
	{{end}}
	</tbody>
	</table>
	</div>
	
	</body>
	</html>
	`)

		if err != nil {
			panic(err)
		}

		results := ResultsPage{
			Term:    *term,
			Results: foundList,
		}

		var exPath string

		if runtime.GOOS == "windows" {
			ex, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}

			exPath = filepath.Dir(ex)
		} else {
			ex, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			exPath = filepath.Dir(ex)
		}

		resultsPath := exPath + string(filepath.Separator) + "_search_results.html"

		file, err := os.OpenFile(resultsPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
		if err != nil {
			log.Panicf("cannot render HTML results: %s", err)
		}

		err = t.Execute(file, results)
		if err != nil {
			log.Panicf("cannot render HTML results: %s", err)
		}

		_ = file.Close()

		fmt.Printf("\nopening results: '%s'", resultsPath)

		if runtime.GOOS == "windows" {
			cmd := exec.Command("cmd", "/c", "start ", resultsPath)

			err = cmd.Run()
			if err != nil {
				log.Panicf("cannot open HTML file: %s", err)
			}
		} else {
			cmd := exec.Command("firefox", "file://"+resultsPath)

			err := cmd.Run()
			if err != nil {
				log.Panicf("cannot open HTML file: %s", err)
			}
		}
	}
}

type ResultsPage struct {
	Term    string
	Results []Find
}

type Find struct {
	Path        string
	AbsPath     string
	Occurrences int
}

func visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	// skip directories
	if f.IsDir() {
		return nil
	}

	ext := filepath.Ext(f.Name())

	var contains bool
	var occurrences int

	switch ext {
	case ".xlsx", ".XLSX":
		contains, occurrences, err = xlsx.Contains(path, *term)
		break
	case ".docx", ".DOCX":
		contains, occurrences, err = docx.Contains(path, *term)
		break
	default:
		return nil
	}

	if err != nil {
		log.Panic(err)
	}

	if contains {
		absPath, err := filepath.Abs(path)

		if err != nil {
			absPath = ""
		}

		foundList = append(foundList, Find{Path: path, AbsPath: absPath, Occurrences: occurrences})
	}

	return nil
}
