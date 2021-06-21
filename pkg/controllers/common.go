package controllers

import (
	"html/template"
	"log"
	"os"
	"fmt"
	"time"
)

func formatDate(args ...interface{}) string {
	for _, item := range args {
		id, ok := item.(time.Time)
		if ok {
			return id.Format("02/01/2006")
		}
	}

	s := fmt.Sprint(args...)
	return s
}

func loadTemplate(pageName string) *template.Template {
	var currentDir, _ = os.Getwd()
	var allFiles = []string{
		currentDir + "/tmpl/common/header.html",
		currentDir + "/tmpl/common/nav.html",
		currentDir + "/tmpl/" + pageName + ".html",
	}

	t, err := template.New("").
		Funcs(template.FuncMap{"formatDate": formatDate}).
		ParseFiles(allFiles...)

	if err != nil {
		log.Panic(err)
	}

	return t
}
