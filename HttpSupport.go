package main

import (
	"net/http"
	"text/template"
)

var tmpl *template.Template

type Data struct {
	Item1 string
	Item2 string
	Item3 string
}

func init() {
	tmpl = template.Must(template.ParseGlob("templates/*.html"))
	//For accessing the assets folder in index.html file . we need to specify that ( IN HTTP ROUTER)
	// fs := http.FileServer(http.Dir("assets"))
	// http.Handle("/assets/", http.StripPrefix("/assets", fs))

	//SAME AS ABOVE FOR MUX ROUTER ( SHIFTED TO MAIN.GO)
	// fs := http.FileServer(http.Dir("assets"))
	// r.PathPrefix("/assets").Handler(http.StripPrefix("/assets", fs))
}

func Dummy_api() string {
	return "HELLO THIS IS YOUR GUILTY CONSCIENCE"
}

func CSS1(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)

}
func CSS2(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "home.html", nil)
}
func DisplayForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "form.html", nil)
}
func DisplayData(w http.ResponseWriter, r *http.Request) {
	recieved_data := Data{
		Item1: r.FormValue("1"),
		Item2: r.FormValue("2"),
		Item3: Dummy_api(),
	}
	tmpl.ExecuteTemplate(w, "display.html", struct {
		Success bool
		Mydata  Data
	}{true, recieved_data})

}
