package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"gitlab.com/bigtallbill/dir_search/pkg/docx"
	"path/filepath"
	"os"
	"gitlab.com/bigtallbill/dir_search/pkg/xlsx"
	"github.com/asticode/go-astilectron"
)

var (
	directory = kingpin.Arg("dir", "the directory to search in").String()
	term      = kingpin.Arg("term", "the search term").String()
	makeHtml  = kingpin.Flag("out", "generates a html file of results and opens it").Short('o').Bool()
	foundList []Find
)

func main() {
	kingpin.Parse()

	// Initialize astilectron
	var a, err = astilectron.New(astilectron.Options{
		AppName: "Doc Search",
		AppIconDefaultPath: "<your .png icon>", // If path is relative, it must be relative to the data directory
		AppIconDarwinPath:  "<your .icns icon>", // Same here
		BaseDirectoryPath: "deps",
	})
	defer a.Close()

	if err != nil {
		log.Print(err)
	}

	// Start astilectron
	a.Start()

	var w, _ = a.NewWindow("http://127.0.0.1:4000", &astilectron.WindowOptions{
		Center: astilectron.PtrBool(true),
		Height: astilectron.PtrInt(600),
		Width:  astilectron.PtrInt(600),
	})
	w.Create()

	// Blocking pattern
	a.Wait()
//
//	http.HandleFunc("/", controllers.Home)
//
//	err := http.ListenAndServe(":9090", nil) // set listen port
//	if err != nil {
//		log.Fatal("ListenAndServe: ", err)
//	}
//
//	err = filepath.Walk(*directory, visit)
//
//	if err != nil {
//		panic(fmt.Sprintf("filepath.Walk() returned %v", err))
//	}
//
//	fmt.Printf("Found %s in %d documents:\n\n", *term, len(foundList))
//	for _, find := range foundList {
//		fmt.Printf("%dx in \"%s\"\n", find.Occurrences, find.Path)
//	}
//
//	if *makeHtml {
//		var t = template.New("home.html") // Create a template.
//		t, err = t.Parse(`
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <title>Search Results for {{.Term}}</title>
//	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.1.1/css/bootstrap.min.css" integrity="sha384-WskhaSGFgHYWDcbwN70/dfYBj47jz9qbsMId/iRN3ewGhXQFZCSftd1LZCfmhktB" crossorigin="anonymous">
//</head>
//<body>
//<div class="container">
//<h2>"{{.Term}}" Results:</h2>
//	<table class="table">
//<thead>
//<tr>
//      <th scope="col">Count</th>
//      <th scope="col">File</th>
//    </tr>
//</thead>
//	<tbody>
//{{range .Results}}
//<tr>
//<td>{{.Occurrences}}</td>
//<td><a href="file:///{{.AbsPath}}">{{.Path}}</a></td>
//</tr>
//{{end}}
//</tbody>
//</table>
//</div>
//
//</body>
//</html>
//`)
//		if err != nil {
//			panic(err)
//		}
//
//		results := ResultsPage{
//			Term:    *term,
//			Results: foundList,
//		}
//
//		var exPath string
//
//		if runtime.GOOS == "windows" {
//			ex, err := filepath.Abs(filepath.Dir(os.Args[0]))
//			if err != nil {
//				log.Fatal(err)
//			}
//			exPath = filepath.Dir(ex)
//		} else {
//			ex, err := os.Getwd()
//			if err != nil {
//				panic(err)
//			}
//			exPath = filepath.Dir(ex)
//		}
//
//		resultsPath := exPath + string(filepath.Separator) + "_search_results.html"
//
//		file, err := os.OpenFile(resultsPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
//		t.Execute(file, results)
//
//		if err != nil {
//			panic(err)
//		}
//
//		file.Close()
//
//		fmt.Printf("\nopening results: '%s'", resultsPath)
//
//		if runtime.GOOS == "windows" {
//			cmd := exec.Command("cmd", "/c", "start ", resultsPath)
//			err = cmd.Run()
//			if err != nil {
//				log.Print(err)
//			}
//		} else {
//			cmd := exec.Command("firefox", "file://"+resultsPath)
//			cmd.Run()
//		}
//	}
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
		log.Panic(err)
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
