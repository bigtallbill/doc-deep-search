package controllers

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t := loadTemplate("home")

	//shopService := services.NewShop(db.New())

	err := t.ExecuteTemplate(
		w,
		"home.html",
		nil,
	)

	if err != nil {
		panic(err)
	}
}
